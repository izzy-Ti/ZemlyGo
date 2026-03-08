package realtime

import (
	"log"
	"sync"
)

type Hub struct {
	mu     sync.RWMutex
	client map[string]*client
}

func NewHub() *Hub {
	return &Hub{
		client: make(map[string]*client),
	}
}

func (h *Hub) Register(client *client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.client[client.userID] = client
	log.Println("client connected: ", client.userID)
}

func (h *Hub) Unregister(userId string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.client, userId)
	log.Println("client disconnected: ", userId)
}
func (h *Hub) SendToUser(userId string, event Event) {
	h.mu.RLock()
	client, ok := h.client[userId]
	h.mu.RUnlock()

	if !ok {
		return
	}
	select {
	case client.Send <- event:
	default:
		log.Println("send channel full for user: ", userId)
	}
}
