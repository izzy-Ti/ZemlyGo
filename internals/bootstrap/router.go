package bootstrap

import (
	"github.com/gin-gonic/gin"
	"github.com/izzy-Ti/ZemlyGo/internals/routes"
)

func NewRoute(app *App, providers *Providers) *gin.Engine {
	if !app.Config.IsDevelopment() {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	// Register all routes and middlewares
	routes.SetupRoutes(r, providers.ToRouteProviders(), app.Config)

	return r
}
