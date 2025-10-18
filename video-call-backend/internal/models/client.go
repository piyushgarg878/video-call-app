package models

import "github.com/gofiber/websocket/v2"

type Client struct {
	ID     string
	Conn   *websocket.Conn
	RoomID string
	Send   chan Message
}