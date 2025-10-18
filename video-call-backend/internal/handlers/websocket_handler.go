package handlers

import (
	"log"

	"github.com/gofiber/websocket/v2"
	"github.com/piyushgarg878/video-call-backend/internal/models"
	"github.com/piyushgarg878/video-call-backend/internal/ws"
)

type WSHandler struct {
	hub *ws.Hub
}

func NewWSHandler(h *ws.Hub) *WSHandler { return &WSHandler{hub: h} }

func (h *WSHandler) Handle(c *websocket.Conn) {
	client := &models.Client{
		ID:     c.Params("userId"), // or generate
		RoomID: c.Query("roomId"),
		Conn:   c,
		Send:   make(chan models.Message, 10),
	}
	h.hub.AddClient(client.RoomID, client)
	defer h.hub.RemoveClient(client.RoomID, client.ID)

	go func() {
		for msg := range client.Send {
			if err := c.WriteJSON(msg); err != nil {
				log.Println("write err:", err)
				return
			}
		}
	}()

	for {
		var msg models.Message
		if err := c.ReadJSON(&msg); err != nil {
			log.Println("read err:", err)
			break
		}

		switch msg.Type {
		case "chat":
			h.hub.Broadcast(client.RoomID, msg, client.ID)
		case "offer", "answer", "ice":
			h.hub.Broadcast(client.RoomID, msg, client.ID)
		default:
			log.Println("unknown message:", msg.Type)
		}
	}
}