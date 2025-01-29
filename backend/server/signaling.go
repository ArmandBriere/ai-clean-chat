package server

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var AllRooms RoomMap

// CreateRoomRequestHandler handles the request to create a new room
func CreateRoomRequestHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	roomID := AllRooms.CreateRoom()

	// Return the roomID as a JSON response
	json.NewEncoder(w).Encode(RoomCreationResponse{RoomID: roomID})
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Message struct {
	Type      string `json:"type"`
	Offer     string `json:"offer,omitempty"`
	Answer    string `json:"answer,omitempty"`
	Candidate struct {
		Candidate        string `json:"candidate"`
		SdpMid           string `json:"sdpMid"`
		SdpMLineIndex    int    `json:"sdpMLineIndex"`
		UsernameFragment string `json:"usernameFragment"`
	} `json:"candidate,omitempty"`
}

type broadcastMsg struct {
	Message map[string]interface{}
	RoomID  string
	UserID  string
	Client  *websocket.Conn
}

var broadcast = make(chan broadcastMsg)

func broadcaster() {
	for {
		msg := <-broadcast
		for id, client := range AllRooms.Map[msg.RoomID] {
			slog.Info("Client ID", "id", id)
			if client.UserID != msg.UserID {

				// Check if the connection is still open
				if err := client.Conn.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(time.Second)); err != nil {
					slog.Error("Client connection is closed", "err", err)
					AllRooms.DeleteFromRoom(msg.RoomID, client.UserID)
					continue
				}

				err := client.Conn.WriteJSON(msg.Message)
				if err != nil {
					slog.Error("An error occur while writing", "err", err)
				}
			}
		}
	}
}
func init() {
	go broadcaster()
}

// JoinRoomRequestHandler handles the request to join a room and listen on the websocket connection
func JoinRoomRequestHandler(w http.ResponseWriter, r *http.Request) {
	roomID, ok := r.URL.Query()["roomID"]
	if !ok {
		slog.Info("RoomID missing in URL Parameters")
		return
	}

	userID, ok := r.URL.Query()["userID"]
	if !ok {
		slog.Info("userID missing in URL Parameters")
		return
	}

	wsConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Error("WebSocket connection upgrade failed", "Error", err)
		return
	}
	defer wsConn.Close()

	slog.Info("New Connection", "roomID", roomID)
	slog.Info("New Connection", "roomID[0]", roomID[0])
	AllRooms.InsertIntoRoom(roomID[0], userID[0], wsConn)

	// This is the main loop that listens for messages from the client
	for {
		var message map[string]interface{}
		err := wsConn.ReadJSON(&message)

		if err != nil {
			slog.Error("Error reading JSON", "err", err)
			return
		}

		broadcastMsg := broadcastMsg{
			Message: message,
			RoomID:  roomID[0],
			UserID:  userID[0],
			Client:  wsConn,
		}

		broadcast <- broadcastMsg
	}
}
