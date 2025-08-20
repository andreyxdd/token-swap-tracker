package ws

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type Client struct {
	upgrader  websocket.Upgrader
	mu        *sync.Mutex
	clients   map[*websocket.Conn]bool
	broadcast chan []byte
}

func New(broadcast chan []byte) *Client {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	return &Client{
		upgrader:  upgrader,
		mu:        &sync.Mutex{},
		clients:   make(map[*websocket.Conn]bool),
		broadcast: broadcast,
	}
}

func (c *Client) Handler(w http.ResponseWriter, r *http.Request) {
	conn, err := c.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error upgrading: %v", err)
		return
	}
	defer func() {
		c.mu.Lock()
		delete(c.clients, conn)
		c.mu.Unlock()
		conn.Close()
	}()

	c.mu.Lock()
	c.clients[conn] = true
	c.mu.Unlock()

	log.Printf("Client connected. Total clients: %d", len(c.clients))
	c.drainIncomingMessages(conn)
	log.Printf("Client disconnected. Total clients: %d", len(c.clients)-1)
}

func (c *Client) HandleBroadcasting() {
	for {
		message := <-c.broadcast
		c.mu.Lock()

		var toRemove []*websocket.Conn
		for client := range c.clients {
			err := client.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				log.Printf("Error writing to client: %v", err)
				toRemove = append(toRemove, client)
			}
		}

		for _, client := range toRemove {
			client.Close()
			delete(c.clients, client)
		}

		c.mu.Unlock()
	}
}

func (c *Client) drainIncomingMessages(conn *websocket.Conn) {
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Unexpected WebSocket error: %v", err)
			}
			break
		}
	}
}
