package webrtcserver

import "github.com/pion/webrtc/v4"

type WebSocketMessage struct {
	UserID           string `json:"userId"`
	Target           string `json:"target"`
	LocalDescription string `json:"localDescription,omitempty"`
}

type User struct {
	UserID         string `json:"userId"`
	PeerConnection *webrtc.PeerConnection
}

type UserMap map[string]*User
