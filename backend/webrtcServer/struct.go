package webrtcserver

type WebSocketMessage struct {
	Type      string `json:"type"`
	SDP       string `json:"sdp"`
	Candidate struct {
		Candidate        string `json:"candidate"`
		SdpMid           string `json:"sdpMid"`
		SdpMLineIndex    int    `json:"sdpMLineIndex"`
		UsernameFragment string `json:"usernameFragment"`
	} `json:"candidate"`
	IsStreaming bool `json:"isStreaming"`
}
