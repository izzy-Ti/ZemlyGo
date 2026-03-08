package bootstrap

import (
	"log"

	"github.com/izzy-Ti/ZemlyGo/configs"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewDatabase(cfg *configs.Config) *gorm.DB {
	db, err := gorm.Open(postgres.Open(cfg.DataBaseURL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatal("failed to connect to database: ", err)
	}
	log.Println("Database connected")
	return db
}
