package bootstrap

import (
	"github.com/izzy-Ti/ZemlyGo/configs"
	"github.com/izzy-Ti/ZemlyGo/internals/infrastructure/neon"
	"github.com/izzy-Ti/ZemlyGo/internals/realtime"
	"gorm.io/gorm"
)

type App struct {
	Config   *configs.Config
	DB       *gorm.DB
	Hub      *realtime.Hub
	NeonAuth *neon.Client
}

func NewApp(cfg *configs.Config, db *gorm.DB) *App {
	neonClient := neon.NewClient(cfg)

	return &App{
		Config:   cfg,
		DB:       db,
		Hub:      realtime.NewHub(),
		NeonAuth: neonClient,
	}
}
