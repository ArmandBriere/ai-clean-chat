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

type broadcastMsg struct {
	Message map[string]interface{}
	RoomID  string
	UserID  string
	Client  *websocket.Conn
}

var broadcast = make(chan broadcastMsg)

// broadcaster is a goroutine that listens for messages from the broadcast channel and sends them to all clients in the room
func broadcaster() {
	for {
		msg := <-broadcast
		for _, client := range AllRooms.Map[msg.RoomID] {
			// Don't send the message back to the sender
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
	roomID := r.URL.Query().Get("roomID")
	if roomID == "" {
		slog.Info("RoomID missing in URL Parameters")
		http.Error(w, "RoomID missing in URL Parameters", http.StatusBadRequest)
		return
	}

	userID := r.URL.Query().Get("userID")
	if userID == "" {
		slog.Info("userID missing in URL Parameters")
		http.Error(w, "userID missing in URL Parameters", http.StatusBadRequest)
		return
	}

	wsConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Error("WebSocket connection upgrade failed", "Error", err)
		return
	}
	defer wsConn.Close()

	AllRooms.InsertIntoRoom(roomID, userID, wsConn)

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
			RoomID:  roomID,
			UserID:  userID,
			Client:  wsConn,
		}

		broadcast <- broadcastMsg
	}
}
