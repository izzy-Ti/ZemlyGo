package main

import (
	"github.com/izzy-Ti/ZemlyGo/configs"
	"github.com/izzy-Ti/ZemlyGo/internals/bootstrap"
	"github.com/izzy-Ti/ZemlyGo/internals/infrastructure/database"
)

func main() {
	// r := gin.Default()

	// r.GET("/ping", func(c *gin.Context) {
	// 	c.JSON(200, gin.H{
	// 		"message": "PONG",
	// 	})
	// })

	// r.Run()

	cfg := configs.Load()

	db := bootstrap.NewDatabase(cfg)
	database.RunMigration(db)

	app := bootstrap.NewApp(cfg, db)
	providers := bootstrap.NewProvider(app)
	router := bootstrap.NewRoute(app, providers)

	router.Run(":" + cfg.Port)
}
