package main

import (
	"fmt"
	"log"
	"net/http"

	server "profanity.com/server"
)

func main() {
	var allRooms server.RoomMap
	allRooms.Init()

	http.HandleFunc("/create", server.CreateRoomRequestHandler)
	http.HandleFunc("/join", server.JoinRoomRequestHandler)

	log.Println("starting Server on  Port 8000")
	fmt.Println(" ")
	err := http.ListenAndServe(":8000", nil)

	if err != nil {
		log.Fatal((err))
	}
}
