package realtime

import "github.com/gorilla/websocket"

type Client struct {
	UserID string
	Role   string
	Conn   *websocket.Conn
	Send   chan Event
}
