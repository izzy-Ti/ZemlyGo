package bootstrap

import "github.com/gin-gonic/gin"

func NewRoute(app *App) *gin.Engine {
	r := gin.Default()

	return r
}
