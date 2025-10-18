package models

type Message struct {
	Type    string      `json:"type"`
	RoomID  string      `json:"roomId"`
	Sender  string      `json:"sender"`
	Target  string      `json:"target,omitempty"`
	Payload interface{} `json:"payload,omitempty"`
}