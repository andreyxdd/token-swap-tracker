package main

import (
	"log"

	"github.com/gorilla/websocket"
)

func main() {
	serverURL := "ws://localhost:8082/ws"
	conn, _, err := websocket.DefaultDialer.Dial(serverURL, nil)
	if err != nil {
		log.Println("Error connecting to WebSocket server: ", err)
	}
	defer conn.Close()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message: ", err)
		}
		log.Println("Received:", string(message))
	}
}
