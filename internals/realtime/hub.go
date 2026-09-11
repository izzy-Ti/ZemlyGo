package realtime

import (
	"log"
	"sync"
)

type Hub struct {
	mu      sync.RWMutex
	Clients map[string]*Client
}

func NewHub() *Hub {
	return &Hub{
		Clients: make(map[string]*Client),
	}
}

func (h *Hub) Register(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	// If client already exists, close previous connection cleanly
	if existing, exists := h.Clients[client.UserID]; exists {
		close(existing.Send)
		_ = existing.Conn.Close()
	}

	h.Clients[client.UserID] = client
	log.Printf("[WebSocket Hub] Client registered: UserID=%s Role=%s\n", client.UserID, client.Role)
}

func (h *Hub) Unregister(userID string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if client, ok := h.Clients[userID]; ok {
		delete(h.Clients, userID)
		close(client.Send)
		log.Printf("[WebSocket Hub] Client unregistered: UserID=%s\n", userID)
	}
}

func (h *Hub) SendToUser(userID string, event Event) {
	h.mu.RLock()
	client, ok := h.Clients[userID]
	h.mu.RUnlock()

	if !ok {
		return
	}

	select {
	case client.Send <- event:
	default:
		log.Printf("[WebSocket Hub] Send channel full for user %s; dropping event\n", userID)
	}
}

func (h *Hub) SendToRole(role string, event Event) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for _, client := range h.Clients {
		if client.Role == role {
			select {
			case client.Send <- event:
			default:
				log.Printf("[WebSocket Hub] Send channel full for user %s; dropping event\n", client.UserID)
			}
		}
	}
}

func (h *Hub) Broadcast(event Event) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for _, client := range h.Clients {
		select {
		case client.Send <- event:
		default:
		}
	}
}

func (h *Hub) IsUserOnline(userID string) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	_, ok := h.Clients[userID]
	return ok
}
