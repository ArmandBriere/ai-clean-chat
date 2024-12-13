package webrtcserver

import (
	"encoding/json"
	"log/slog"

	"github.com/gorilla/websocket"
	"github.com/pion/webrtc/v4"
)

// parseOfferMessage parses the offer message
func parseOfferMessage(msg WebSocketMessage, peerConnection *webrtc.PeerConnection, conn *websocket.Conn) {
	slog.Info("Offer message received")
	sdp := msg.SDP
	offer := webrtc.SessionDescription{Type: webrtc.SDPTypeOffer, SDP: sdp}
	if err := peerConnection.SetRemoteDescription(offer); err != nil {
		slog.Error("Error setting remote description", "error", err)
		return
	}

	answer, err := peerConnection.CreateAnswer(nil)
	if err != nil {
		slog.Error("Error creating the answer", "error", err)
		return
	}

	if err := peerConnection.SetLocalDescription(answer); err != nil {
		slog.Error("Error setting local description", "error", err)
		return
	}

	if err := conn.WriteJSON(WebSocketMessage{Type: "answer", SDP: answer.SDP}); err != nil {
		slog.Error("Error writing answer", "error", err)
		return
	}
}

// parseIceCandidateMessage parses the ICE candidate message
func parseIceCandidateMessage(msg WebSocketMessage, peerConnection *webrtc.PeerConnection) {
	candidateData := msg.Candidate
	if candidateData == "" {
		slog.Error("No ICE candidate")
		return
	}

	candidateBytes, err := json.Marshal(candidateData)
	if err != nil {
		slog.Error("Failed to marshal candidate data", "error", err)
		return
	}

	var candidate webrtc.ICECandidateInit
	if err := json.Unmarshal(candidateBytes, &candidate); err != nil {
		slog.Error("Failed to unmarshal candidate:", "error", err)
		return
	}

	if err := peerConnection.AddICECandidate(candidate); err != nil {
		slog.Error("Error adding ICE candidate:", "error", err)
		return
	}
}

// parseStreamingMessage parses the streaming message
func parseStreamingMessage(isStreaming *bool, msg WebSocketMessage) {
	slog.Info("Streaming message received")
	*isStreaming = msg.IsStreaming
	if *isStreaming {
		slog.Info("Starting streaming")
	} else {
		slog.Info("Stopping streaming")
	}
}
