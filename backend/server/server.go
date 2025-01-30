package server

import (
	"log/slog"
	"math/rand"
	"time"

	"github.com/gorilla/websocket"
)

func (r *RoomMap) Init() {
	r.Map = make(map[string][]Participant)
}

// Get returns the list of participants in a room
func (r *RoomMap) Get(roomID string) []Participant {
	r.Mutex.RLock()
	defer r.Mutex.RUnlock()

	return r.Map[roomID]
}

// CreateRoom creates a new room and returns the roomID
func (r *RoomMap) CreateRoom() string {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()
	rand.New(rand.NewSource(time.Now().UnixNano()))

	// Generate a random roomID following the pattern: XXX-XXXX-XXX
	roomID := randStringRunes(3) + "-" + randStringRunes(4) + "-" + randStringRunes(3)

	r.Map[roomID] = []Participant{}

	return roomID
}

// randStringRunes generates a random string of length n
func randStringRunes(n int) string {
	var letters = []rune("abcdefghijklmnpqrstuvwxyz")

	if n > len(letters) {
		slog.Error("n is greater than the length of the letters array")
	}

	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// InsertIntoRoom inserts a new participant into the room
func (r *RoomMap) InsertIntoRoom(roomID string, userID string, conn *websocket.Conn) {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	p := Participant{userID, conn}

	slog.Info("Inserting into Room", "roomID", roomID)
	r.Map[roomID] = append(r.Map[roomID], p)
}

// DeleteFromRoom deletes a participant from the room
func (r *RoomMap) DeleteFromRoom(roomID string, userID string) {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	for i, p := range r.Map[roomID] {
		if p.UserID == userID {
			slog.Info("Deleting from Room", "roomID", roomID, "userID", userID)
			r.Map[roomID] = append(r.Map[roomID][:i], r.Map[roomID][i+1:]...)
			return
		}
	}

	// Delete the room if there are no participants left
	if len(r.Map[roomID]) == 0 {
		r.DeleteRoom(roomID)
	}
}

// DeleteRoom deletes a room from the map
func (r *RoomMap) DeleteRoom(roomID string) {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	slog.Info("Deleting Room", "roomID", roomID)
	delete(r.Map, roomID)
}
