package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	server "profanity.com/server"
	webrtcServer "profanity.com/webrtcServer"
)

// health is a simple health check handler
func health(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Health Check")
}

func main() {
	server.AllRooms.Init()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/health", health)
	http.HandleFunc("/create", server.CreateRoomRequestHandler)
	http.HandleFunc("/join", server.JoinRoomRequestHandler)

	// Add the WebRTC handle for transcription
	webrtcServer.AddWebRTCHandle()

	log.Println("Starting server on port " + port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal((err))
	}
}
