package webrtcserver

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pion/webrtc/v4"
	"github.com/pion/webrtc/v4/pkg/media/oggwriter"
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

	peerConnection.OnTrack(func(track *webrtc.TrackRemote, receiver *webrtc.RTPReceiver) {
		codecName := track.Codec().MimeType
		slog.Info("Got track, codec", "codecName", codecName)

		if codecName == webrtc.MimeTypeOpus {
			slog.Info("Track has started")
			go func() {
				filename := fmt.Sprintf("audio_%d.ogg", time.Now().UnixNano())
				f, err := os.Create(filename)
				if err != nil {
					slog.Error("Error creating file", "Error", err)
					return
				}
				defer f.Close()

				// Create an OggWriter
				writer, err := oggwriter.New(filename, 48000, 2)
				if err != nil {
					slog.Error("Error creating oggwriter", "Error", err)
					return
				}
				defer writer.Close()

				for {
					rtpPacket, _, err := track.ReadRTP()
					if err != nil {
						slog.Error("rtpPacket setup", "Error", err)
						return
					}
					if err := writer.WriteRTP(rtpPacket); err != nil {
						slog.Error("write RTP", "Error", err)
						return
					}
				}
			}()
		}
	})

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			slog.Error("Read message error", "error", err)
			break
		}

		var msg map[string]interface{}
		if err := json.Unmarshal(message, &msg); err != nil {
			slog.Error("JSON marshal error", "Error", err)
			continue
		}

		switch msg["type"] {
		case "offer":
			sdp := msg["sdp"].(string)
			offer := webrtc.SessionDescription{Type: webrtc.SDPTypeOffer, SDP: sdp}
			if err := peerConnection.SetRemoteDescription(offer); err != nil {
				slog.Error("Error setting remote description", "error", err)
				continue
			}

			answer, err := peerConnection.CreateAnswer(nil)
			if err != nil {
				slog.Error("Error creating the answer", "error", err)
				continue
			}

			if err := peerConnection.SetLocalDescription(answer); err != nil {
				slog.Error("Error setting local description", "error", err)
				continue
			}

			if err := conn.WriteJSON(map[string]interface{}{"type": "answer", "sdp": answer.SDP}); err != nil {
				slog.Error("Error writing answer", "error", err)
				continue
			}
		case "iceCandidate":
			candidateData, ok := msg["candidate"].(map[string]interface{})
			if !ok {
				slog.Error("Invalid ICE candidate format", "candidate", msg["candidate"])
				continue
			}

			candidateBytes, err := json.Marshal(candidateData)
			if err != nil {
				slog.Error("Failed to marshal candidate data", "error", err)
				continue
			}

			var candidate webrtc.ICECandidateInit
			if err := json.Unmarshal(candidateBytes, &candidate); err != nil {
				slog.Error("Failed to unmarshal candidate:", "error", err)
				continue
			}

			if err := peerConnection.AddICECandidate(candidate); err != nil {
				slog.Error("Error adding ICE candidate:", "error", err)
				continue
			}
		}
	}
}
