package server

import (
	"sync"

	"github.com/gorilla/websocket"
)

type Participant struct {
	UserID string
	Conn   *websocket.Conn
}

type RoomMap struct {
	Mutex sync.RWMutex
	Map   map[string][]Participant
}
