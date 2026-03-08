package bootstrap

import (
	"github.com/izzy-Ti/ZemlyGo/configs"
	"gorm.io/gorm"
)

type App struct {
	config *configs.Config
	DB     *gorm.DB
}

func NewApp(cfg *configs.Config, db *gorm.DB) *App {
	return &App{
		config: cfg,
		DB:     db,
	}
}
