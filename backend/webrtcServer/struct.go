package webrtcserver

type WebSocketMessage struct {
	Type      string `json:"type"`
	SDP       string `json:"sdp,omitempty"`
	Candidate struct {
		Candidate        string `json:"candidate"`
		SdpMid           string `json:"sdpMid"`
		SdpMLineIndex    int    `json:"sdpMLineIndex"`
		UsernameFragment string `json:"usernameFragment"`
	} `json:"candidate,omitempty"`
	IsStreaming bool `json:"isStreaming,omitempty"`
}

type WebSocketTranscription struct {
	Type string `json:"type"`
	Text string `json:"text"`
}
