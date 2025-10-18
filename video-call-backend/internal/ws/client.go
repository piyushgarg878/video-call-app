package ws

import (
	"encoding/json"
	"log"

	"github.com/gofiber/websocket/v2"
)

type Client struct {
	ID     string
	Conn   *websocket.Conn
	Send   chan []byte
	Hub    *Hub
	RoomID string
}

type Message struct {
	Type string          `json:"type"`
	Room string          `json:"room,omitempty"`
	Data json.RawMessage `json:"data,omitempty"`
}

func (c *Client) ReadPump() {
	defer func() {
		c.Hub.RemoveClientFromAllRooms(c)
		c.Conn.Close()
	}()

	for {
		_, msg, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println("read error:", err)
			break
		}

		var m Message
		if err := json.Unmarshal(msg, &m); err != nil {
			log.Println("unmarshal error:", err)
			continue
		}

		switch m.Type {
		case "join-room":
			room := c.Hub.GetRoom(m.Room)
			c.RoomID = m.Room
			room.AddClient(c)
			log.Println("Client joined room:", m.Room)

		case "offer", "answer", "candidate", "leave", "chat":
			room := c.Hub.GetRoom(c.RoomID)
			room.Broadcast(c, msg)
		}
	}
}

func (c *Client) WritePump() {
	for msg := range c.Send {
		if err := c.Conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			log.Println("write error:", err)
			break
		}
	}
}