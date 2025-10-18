package ws

import "sync"

type Hub struct {
	Rooms map[string]*Room
	mu    sync.Mutex
}

func NewHub() *Hub {
	return &Hub{
		Rooms: make(map[string]*Room),
	}
}

func (h *Hub) GetRoom(id string) *Room {
	h.mu.Lock()
	defer h.mu.Unlock()

	room, exists := h.Rooms[id]
	if !exists {
		room = NewRoom(id)
		h.Rooms[id] = room
	}
	return room
}

func (h *Hub) RemoveClientFromAllRooms(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	for _, room := range h.Rooms {
		room.RemoveClient(client)
	}
}