//go:build !profanity

package webrtcserver

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/pion/webrtc/v4"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true }, // For development, REMOVE IN PRODUCTION
}

var peerConnections = make(map[*websocket.Conn]*webrtc.PeerConnection)
var mu sync.Mutex

// AddWebRTCHandle starts the WebRTC server
func AddWebRTCHandle() {
	http.HandleFunc("/ws", handleWebSocket)
}

// handleWebSocket handles incoming WebRTC connections
func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	wsConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Error("WebSocket connection upgrade failed", "Error", err)
		return
	}
	defer wsConn.Close()

	// Register the MediaEngine
	mediaEngine := webrtc.MediaEngine{}
	if err := mediaEngine.RegisterDefaultCodecs(); err != nil {
		slog.Error("Codec register failed", "Error", err)
		return
	}

	api := webrtc.NewAPI(webrtc.WithMediaEngine(&mediaEngine))

	// Create a new RTCPeerConnection
	peerConnection, err := api.NewPeerConnection(webrtc.Configuration{})
	if err != nil {
		slog.Error("New peer connection failed", "Error", err)
		return
	}
	mu.Lock()
	peerConnections[wsConn] = peerConnection
	mu.Unlock()

	defer func() {
		mu.Lock()
		delete(peerConnections, wsConn)
		mu.Unlock()
		peerConnection.Close()
	}()

	// Listen for ICE candidates and write them to the WebSocket
	peerConnection.OnICECandidate(func(i *webrtc.ICECandidate) {
		if i == nil {
			return
		}

		candidate, err := json.Marshal(i.ToJSON())
		if err != nil {
			slog.Error("JSON marshal error", "Error", err)
			return
		}

		mu.Lock()
		defer mu.Unlock()
		if err := wsConn.WriteJSON(map[string]interface{}{"type": "iceCandidate", "candidate": string(candidate)}); err != nil {
			slog.Error("Writing iceCandidate failed", "Error", err)
			return
		}
	})

	// Streaming flag to start/stop audio transcription
	isStreaming := false

	// Done signal that stops the transcription and delete resources
	ctx, cancel := context.WithCancel(context.Background())

	// Handle incoming audio
	peerConnection.OnTrack(func(track *webrtc.TrackRemote, receiver *webrtc.RTPReceiver) {
		codecName := track.Codec().MimeType
		slog.Info("Got track, codec", "codecName", codecName)

		if codecName == webrtc.MimeTypeOpus {
			slog.Info("Track has started")

			go handleAudioStream(ctx, track, &isStreaming, wsConn, &mu)
		}
	})

	for {
		_, message, err := wsConn.ReadMessage()
		if err != nil {
			slog.Error("Read message error", "error", err)
			slog.Info("Cancel the context")
			cancel()
			return
		}

		var msg WebSocketMessage
		if err := json.Unmarshal(message, &msg); err != nil {
			slog.Error("JSON marshal error", "message", message, "error", err)
			continue
		}

		switch msg.Type {
		case "offer":
			parseOfferMessage(msg, peerConnection, wsConn, &mu)
		case "iceCandidate":
			parseIceCandidateMessage(msg, peerConnection)
		case "streaming":
			parseStreamingMessage(&isStreaming, msg)
			message, err := json.Marshal(WebSocketMessage{Type: "streaming", IsStreaming: isStreaming})
			if err != nil {
				slog.Error("JSON marshal error", "Error", err)
				continue
			}
			mu.Lock()
			wsConn.WriteMessage(websocket.TextMessage, message)
			mu.Unlock()
		}
	}
}
