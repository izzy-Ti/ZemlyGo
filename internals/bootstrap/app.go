package bootstrap

import (
	"github.com/izzy-Ti/ZemlyGo/configs"
	"github.com/izzy-Ti/ZemlyGo/internals/realtime"
	"gorm.io/gorm"
)

type App struct {
	config *configs.Config
	DB     *gorm.DB
	Hub    *realtime.Hub
}

func NewApp(cfg *configs.Config, db *gorm.DB) *App {
	return &App{
		config: cfg,
		DB:     db,
		Hub:    realtime.NewHub(),
	}
}
