package handlers

import (
	"github.com/gofiber/websocket/v2"
	"github.com/piyushgarg878/video-call-backend/internal/ws"
)

type WSHandler struct {
	hub *ws.Hub
}

func NewWSHandler(h *ws.Hub) *WSHandler {
	return &WSHandler{hub: h}
}

func (h *WSHandler) Handle(c *websocket.Conn) {
	client := &ws.Client{
		ID:   c.RemoteAddr().String(),
		Conn: c,
		Send: make(chan []byte, 256),
		Hub:  h.hub,
	}

	go client.WritePump()
	client.ReadPump()
}