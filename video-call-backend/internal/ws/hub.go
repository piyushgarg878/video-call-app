package ws

import (
	"sync"

	"github.com/piyushgarg878/video-call-backend/internal/models"
)

type Hub struct {
	rooms map[string]map[string]*models.Client
	mu    sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		rooms: make(map[string]map[string]*models.Client),
	}
}

func (h *Hub) AddClient(roomID string, c *models.Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.rooms[roomID] == nil {
		h.rooms[roomID] = make(map[string]*models.Client)
	}
	h.rooms[roomID][c.ID] = c
}

func (h *Hub) RemoveClient(roomID, clientID string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if room, ok := h.rooms[roomID]; ok {
		delete(room, clientID)
		if len(room) == 0 {
			delete(h.rooms, roomID)
		}
	}
}

func (h *Hub) Broadcast(roomID string, msg models.Message, exclude string) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	if room, ok := h.rooms[roomID]; ok {
		for id, c := range room {
			if id == exclude { continue }
			select {
			case c.Send <- msg:
			default:
			}
		}
	}
}