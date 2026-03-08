package realtime

import (
	"log"
	"sync"
)

type Hub struct {
	mu     sync.RWMutex
	Client map[string]*Client
}

func NewHub() *Hub {
	return &Hub{
		Client: make(map[string]*Client),
	}
}

func (h *Hub) Register(Client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.Client[Client.UserID] = Client
	log.Println("Client connected: ", Client.UserID)
}

func (h *Hub) Unregister(UserId string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.Client, UserId)
	log.Println("Client disconnected: ", UserId)
}
func (h *Hub) SendToUser(UserId string, event Event) {
	h.mu.RLock()
	Client, ok := h.Client[UserId]
	h.mu.RUnlock()

	if !ok {
		return
	}
	select {
	case Client.Send <- event:
	default:
		log.Println("send channel full for user: ", UserId)
	}
}
