package ws

import "sync"

type Room struct {
	ID      string
	Clients map[*Client]bool
	mu      sync.Mutex
}

func NewRoom(id string) *Room {
	return &Room{
		ID:      id,
		Clients: make(map[*Client]bool),
	}
}

func (r *Room) AddClient(client *Client) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.Clients[client] = true
}

func (r *Room) RemoveClient(client *Client) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.Clients, client)
}

func (r *Room) Broadcast(sender *Client, message []byte) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for client := range r.Clients {
		if client != sender {
			client.Send <- message
		}
	}
}