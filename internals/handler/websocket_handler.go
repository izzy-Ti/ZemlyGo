package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/izzy-Ti/ZemlyGo/internals/realtime"
)

type WebsocketHandler struct {
	hub *realtime.Hub
}

func NewWebSocketHandler(hub *realtime.Hub) *WebsocketHandler {
	return &WebsocketHandler{hub: hub}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *WebsocketHandler) Handle(c *gin.Context) {
	userId := c.Query("user_id")
	role := c.Query("role")

	if userId == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "user_id is required",
		})
	}
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("upgrade error: ", err)
		return
	}
	client := &realtime.Client{
		UserID: userId,
		Role:   role,
		Conn:   conn,
		Send:   make(chan realtime.Event),
	}

	h.hub.Register(client)

	go h.writePump(client)
	h.ReadPump(client)
}

func (h *WebsocketHandler) ReadPump(client *realtime.Client) {
	defer func() {
		h.hub.Unregister(client.UserID)
		client.Conn.Close()
	}()
	for {
		var msg map[string]interface{}
		if err := client.Conn.ReadJSON(&msg); err != nil {
			log.Println("read error: ", err)
			break
		}
		log.Println("message from", client.UserID, msg)
	}
}

func (h *WebsocketHandler) writePump(client *realtime.Client) {
	defer client.Conn.Close()

	for event := range client.Send {
		if err := client.Conn.WriteJSON(event); err != nil {
			log.Println("write error: ", err)
			break
		}
	}
}
