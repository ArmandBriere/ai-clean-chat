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
	Type           string  `json:"type"`
	Text           string  `json:"text"`
	Uuid           string  `json:"uuid"`
	ProfanityScore float64 `json:"profanity_score"`
}

// Profanity structs
type PostData struct {
	Text string `json:"text"`
}

type PostResponse struct {
	ProfanityScore float64 `json:"profanity_score"`
}

type LLMAnalysis struct {
	Type        string `json:"type"`
	LLMMessage  string `json:"llm_analysis"`
	UserMessage string `json:"user_message"`
	Timestamp   string `json:"timestamp"`
}
