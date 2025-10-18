package ws

import (
	"encoding/json"
	"log"

	"github.com/gofiber/websocket/v2"
)

// Client represents a connected user over WebSocket
type Client struct {
	ID     string
	Conn   *websocket.Conn
	Send   chan []byte
	Hub    *Hub
	RoomID string
}

// Message is the structure for signaling messages
type Message struct {
	Type   string          `json:"type"`
	Room   string          `json:"room,omitempty"`
	Target string          `json:"target,omitempty"`
	Data   json.RawMessage `json:"data,omitempty"`
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

			// Send all existing peers to new client
			for client := range room.Clients {
				if client != c {
					data, _ := json.Marshal(map[string]string{"id": client.ID})
					c.Send <- wrap("existing-peer", data)
				}
			}

			// Notify all others about new peer
			data, _ := json.Marshal(map[string]string{"id": c.ID})
			room.Broadcast(c, wrap("new-peer", data))

		case "offer", "answer", "candidate":
    		log.Printf("Forwarding %s from %s to %s", m.Type, c.ID, m.Target)
    		room := c.Hub.GetRoom(c.RoomID)
    		for client := range room.Clients {
        		if client.ID == m.Target {
            		client.Send <- msg
        		}
    		}

		case "leave":
			room := c.Hub.GetRoom(c.RoomID)
			room.RemoveClient(c)
		}
	}
}

func (c *Client) WritePump() {
	defer c.Conn.Close()
	for msg := range c.Send {
		if err := c.Conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			break
		}
	}
}

func wrap(t string, data []byte) []byte {
	out, _ := json.Marshal(map[string]interface{}{
		"type": t,
		"data": json.RawMessage(data),
	})
	return out
}