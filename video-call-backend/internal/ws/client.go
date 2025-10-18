package ws

import (
	"encoding/json"
	"log"

	"github.com/gofiber/websocket/v2"
	"github.com/pion/webrtc/v3"
)

// Client represents a connected user over WebSocket
type Client struct {
	ID     string
	Conn   *websocket.Conn
	Send   chan []byte
	Hub    *Hub
	RoomID string
	Peer   *webrtc.PeerConnection
}

// Message is the structure for signaling messages
type Message struct {
	Type string          `json:"type"`
	Room string          `json:"room,omitempty"`
	Data json.RawMessage `json:"data,omitempty"`
}

// ReadPump handles incoming WebSocket messages
func (c *Client) ReadPump() {
	defer func() {
		c.Hub.RemoveClientFromAllRooms(c)
		c.Conn.Close()
		if c.Peer != nil {
			c.Peer.Close()
		}
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

		// --- ROOM JOIN ---
		case "join-room":
			room := c.Hub.GetRoom(m.Room)
			c.RoomID = m.Room
			room.AddClient(c)
			log.Println("Client joined room:", m.Room)

		// --- CHAT MESSAGE ---
		case "chat":
			room := c.Hub.GetRoom(c.RoomID)
			room.Broadcast(c, msg)

		// --- WEBRTC OFFER ---
		case "offer":
			var offer webrtc.SessionDescription
			if err := json.Unmarshal(m.Data, &offer); err != nil {
				log.Println("invalid offer:", err)
				continue
			}

			pc, err := NewPeerConnection()
			if err != nil {
				log.Println("peer conn error:", err)
				continue
			}
			c.Peer = pc

			// Remote description from browser
			if err := pc.SetRemoteDescription(offer); err != nil {
				log.Println("set remote desc error:", err)
				continue
			}

			// Create answer
			answer, err := pc.CreateAnswer(nil)
			if err != nil {
				log.Println("create answer error:", err)
				continue
			}

			if err := pc.SetLocalDescription(answer); err != nil {
				log.Println("set local desc error:", err)
				continue
			}

			// Send answer back to client
			answerJSON, _ := json.Marshal(struct {
				Type string                      `json:"type"`
				Data webrtc.SessionDescription   `json:"data"`
			}{
				Type: "answer",
				Data: answer,
			})
			c.Send <- answerJSON

			// Handle ICE candidates from backend
			pc.OnICECandidate(func(i *webrtc.ICECandidate) {
				if i == nil {
					return
				}
				candidateJSON, _ := json.Marshal(struct {
					Type string                    `json:"type"`
					Data webrtc.ICECandidateInit   `json:"data"`
				}{
					Type: "candidate",
					Data: i.ToJSON(),
				})
				c.Send <- candidateJSON
			})

		// --- WEBRTC CANDIDATE ---
		case "candidate":
			var candidate webrtc.ICECandidateInit
			if err := json.Unmarshal(m.Data, &candidate); err != nil {
				log.Println("invalid candidate:", err)
				continue
			}
			if c.Peer != nil {
				if err := c.Peer.AddICECandidate(candidate); err != nil {
					log.Println("add candidate error:", err)
				}
			}

		// --- LEAVE ROOM ---
		case "leave":
			room := c.Hub.GetRoom(c.RoomID)
			room.RemoveClient(c)
			log.Println("Client left room:", c.RoomID)
		}
	}
}

// WritePump sends messages from the server to the WebSocket client
func (c *Client) WritePump() {
	defer c.Conn.Close()

	for msg := range c.Send {
		if err := c.Conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			log.Println("write error:", err)
			break
		}
	}
}