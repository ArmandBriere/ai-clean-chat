package webrtcserver

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/pion/webrtc/v4"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true }, // For development, REMOVE IN PRODUCTION
}

var peerConnections = make(map[*websocket.Conn]*webrtc.PeerConnection)

// AddWebRTCHandle starts the WebRTC server
func AddWebRTCHandle() {
	http.HandleFunc("/ws", handleWebSocket)
}

// handleWebSocket handles incoming WebRTC connections
func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Error("WebSocket connection upgrade failed", "Error", err)
		return
	}
	defer conn.Close()

	mediaEngine := webrtc.MediaEngine{}
	if err := mediaEngine.RegisterDefaultCodecs(); err != nil {
		slog.Error("Codec register failed", "Error", err)
		return
	}

	api := webrtc.NewAPI(webrtc.WithMediaEngine(&mediaEngine))

	peerConnection, err := api.NewPeerConnection(webrtc.Configuration{})
	if err != nil {
		slog.Error("New peer connection failed", "Error", err)
		return
	}
	peerConnections[conn] = peerConnection
	defer delete(peerConnections, conn)
	defer peerConnection.Close()

	audioTrack, err := webrtc.NewTrackLocalStaticRTP(webrtc.RTPCodecCapability{MimeType: webrtc.MimeTypeOpus}, "audio", "pion")
	if err != nil {
		slog.Error("New audio track failed", "Error", err)
		return
	}

	rtpSender, err := peerConnection.AddTrack(audioTrack)
	if err != nil {
		slog.Error("Adding track failed", "Error", err)
		return
	}

	go func() {
		rtpBuf := make([]byte, 1500)
		for {
			_, _, readErr := rtpSender.Read(rtpBuf)
			if readErr != nil {
				return
			}
		}
	}()

	peerConnection.OnICECandidate(func(i *webrtc.ICECandidate) {
		if i == nil {
			return
		}

		candidate, err := json.Marshal(i.ToJSON())
		if err != nil {
			slog.Error("JSON marshal error", "Error", err)
			return
		}

		if err := conn.WriteJSON(map[string]interface{}{"type": "iceCandidate", "candidate": string(candidate)}); err != nil {
			slog.Error("Writing iceCandidate failed", "Error", err)
			return
		}
	})

	isStreaming := false

	peerConnection.OnTrack(func(track *webrtc.TrackRemote, receiver *webrtc.RTPReceiver) {
		codecName := track.Codec().MimeType
		slog.Info("Got track, codec", "codecName", codecName)

		if codecName == webrtc.MimeTypeOpus {
			slog.Info("Track has started")
			go handleAudioStream(track, &isStreaming)
		}
	})

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			slog.Error("Read message error", "error", err)
			break
		}

		var msg WebSocketMessage
		if err := json.Unmarshal(message, &msg); err != nil {
			slog.Error("JSON marshal error", "message", message, "error", err)
			continue
		}

		switch msg.Type {
		case "offer":
			parseOfferMessage(msg, peerConnection, conn)
		case "iceCandidate":
			parseIceCandidateMessage(msg, peerConnection)
		case "streaming":
			parseStreamingMessage(&isStreaming, msg)
		}
	}
}
