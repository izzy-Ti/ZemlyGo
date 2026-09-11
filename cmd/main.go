package main

import (
	"log"

	"github.com/izzy-Ti/ZemlyGo/configs"
	"github.com/izzy-Ti/ZemlyGo/internals/bootstrap"
	"github.com/izzy-Ti/ZemlyGo/internals/infrastructure/database"
)

func main() {
	log.Println("Starting ZemlyGo Ride-Hailing Service...")

	// 1. Load configuration
	cfg := configs.Load()

	// 2. Initialize database
	db := bootstrap.NewDatabase(cfg)

	// 3. Run schema migrations
	if err := database.RunMigration(db); err != nil {
		log.Printf("Warning: Schema migration encountered errors: %v\n", err)
	}

	// 4. Initialize application core
	app := bootstrap.NewApp(cfg, db)

	// 5. Wire providers (repositories, services, handlers, middleware)
	providers := bootstrap.NewProvider(app)

	// 6. Setup router and routes
	router := bootstrap.NewRoute(app, providers)

	log.Printf("ZemlyGo server listening on port %s in %s mode\n", cfg.Port, cfg.Env)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Server failed to run: %v\n", err)
	}
}
