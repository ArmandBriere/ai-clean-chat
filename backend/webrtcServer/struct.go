package webrtcserver

type WebSocketMessage struct {
	Type        string `json:"type"`
	SDP         string `json:"sdp"`
	Candidate   string `json:"candidate"`
	IsStreaming bool   `json:"isStreaming"`
}
