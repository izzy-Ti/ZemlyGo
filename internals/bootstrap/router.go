package bootstrap

import "github.com/gin-gonic/gin"

func NewRoute(app *App, providers *Providers) *gin.Engine {
	r := gin.Default()

	r.GET("/ws", providers.WebsocketHandler.Handle)

	return r
}
