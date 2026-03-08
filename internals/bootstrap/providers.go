package bootstrap

import "github.com/izzy-Ti/ZemlyGo/internals/handler"

type Providers struct {
	WebsocketHandler *handler.WebsocketHandler
}

func NewProvider(app *App) *Providers {
	WsHandler := handler.NewWebSocketHandler(app.Hub)

	return &Providers{
		WebsocketHandler: WsHandler,
	}
}
