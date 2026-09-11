package handler

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/izzy-Ti/ZemlyGo/internals/realtime"
	"github.com/izzy-Ti/ZemlyGo/internals/utils"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

type WebsocketHandler struct {
	hub       *realtime.Hub
	jwtSecret string
}

func NewWebSocketHandler(hub *realtime.Hub, jwtSecret string) *WebsocketHandler {
	return &WebsocketHandler{
		hub:       hub,
		jwtSecret: jwtSecret,
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Configurable or permissive in development
	},
}

func (h *WebsocketHandler) Handle(c *gin.Context) {
	tokenStr := c.Query("token")
	userID := c.Query("user_id")
	role := c.Query("role")

	// If token is provided, extract authenticated user identity
	if tokenStr != "" {
		if claims, err := utils.ValidateJWT(h.jwtSecret, tokenStr); err == nil {
			userID = strconv.FormatUint(uint64(claims.UserID), 10)
			role = claims.Role
		}
	}

	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "user_id or valid token is required for WebSocket connection",
		})
		return
	}

	if role == "" {
		role = "rider"
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("[WebSocket] upgrade error:", err)
		return
	}

	client := &realtime.Client{
		UserID: userID,
		Role:   role,
		Conn:   conn,
		Send:   make(chan realtime.Event, 64),
	}

	h.hub.Register(client)

	go h.writePump(client)
	go h.readPump(client)
}

func (h *WebsocketHandler) readPump(client *realtime.Client) {
	defer func() {
		h.hub.Unregister(client.UserID)
		_ = client.Conn.Close()
	}()

	client.Conn.SetReadLimit(maxMessageSize)
	_ = client.Conn.SetReadDeadline(time.Now().Add(pongWait))
	client.Conn.SetPongHandler(func(string) error {
		_ = client.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		var msg realtime.Event
		if err := client.Conn.ReadJSON(&msg); err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("[WebSocket] read error for user %s: %v\n", client.UserID, err)
			}
			break
		}

		// Handle client-sent events (e.g. heartbeat ping or location)
		if msg.Type == "ping" {
			client.Send <- realtime.Event{
				Type: "pong",
				Data: map[string]interface{}{"time": time.Now().Unix()},
			}
		}
	}
}

func (h *WebsocketHandler) writePump(client *realtime.Client) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		_ = client.Conn.Close()
	}()

	for {
		select {
		case event, ok := <-client.Send:
			_ = client.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel
				_ = client.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := client.Conn.WriteJSON(event); err != nil {
				log.Printf("[WebSocket] write error for user %s: %v\n", client.UserID, err)
				return
			}

		case <-ticker.C:
			_ = client.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := client.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
