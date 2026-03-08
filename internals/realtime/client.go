package realtime

import "golang.org/x/net/websocket"

type client struct {
	userID string
	Role   string
	Conn   *websocket.Conn
	Send   chan Event
}
