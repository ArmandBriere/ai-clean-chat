package server

import (
	"log"
	"log/slog"
	"math/rand"
	"time"

	"github.com/gorilla/websocket"
)

func (r *RoomMap) Init() {
	r.Map = make(map[string][]Participant)
}

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

func (r *RoomMap) InsertIntoRoom(roomID string, host bool, conn *websocket.Conn) {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	p := Participant{host, conn}

	log.Println("Inserting into Room with RoomID: ", roomID)
	r.Map[roomID] = append(r.Map[roomID], p)
}

func (r *RoomMap) DeleteRoom(roomID string) {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	delete(r.Map, roomID)

}
