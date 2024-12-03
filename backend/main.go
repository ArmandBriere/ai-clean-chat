package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	server "profanity.com/server"
)

func health(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Health Check")
}

func main() {
	var allRooms server.RoomMap

	port := os.Getenv("PORT")

	allRooms.Init()

	http.HandleFunc("/health", health)
	http.HandleFunc("/create", server.CreateRoomRequestHandler)
	http.HandleFunc("/join", server.JoinRoomRequestHandler)

	log.Println("starting Server on  Port " + port)
	fmt.Println(" ")
	err := http.ListenAndServe(":"+port, nil)

	if err != nil {
		log.Fatal((err))
	}
}
