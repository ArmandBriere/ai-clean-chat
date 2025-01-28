package server

import (
	"encoding/json"
	"log"
	"log/slog"
	"net/http"

	"github.com/gorilla/websocket"
)

var AllRooms RoomMap

func CreateRoomRequestHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	roomID := AllRooms.CreateRoom()
	type resp struct {
		RoomID string `json:"room_id"`
	}

	log.Println(AllRooms.Map)
	json.NewEncoder(w).Encode(resp{RoomID: roomID})
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
		for _, client := range AllRooms.Map[msg.RoomID] {
			if client.UserID != msg.UserID {
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

	ws, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		slog.Error("Web Socket Upgrade Error", "err", err)
		return
	}

	slog.Info("New Connection", "roomID", roomID)
	slog.Info("New Connection", "roomID[0]", roomID[0])
	AllRooms.InsertIntoRoom(roomID[0], userID[0], ws)

	for {
		var message map[string]interface{}
		err := ws.ReadJSON(&message)
		if err != nil {
			slog.Error("Error reading JSON", "err", err)
			return
		}

		broadcastMsg := broadcastMsg{
			Message: message,
			RoomID:  roomID[0],
			UserID:  userID[0],
			Client:  ws,
		}

		// slog.Info("New message", "msg", broadcastMsg.Message)

		broadcast <- broadcastMsg
	}
}
