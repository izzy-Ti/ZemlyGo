package database

import (
	"log"

	"github.com/izzy-Ti/ZemlyGo/internals/domain"
	"gorm.io/gorm"
)

func RunMigration(db *gorm.DB) error {
	err := db.AutoMigrate(
		&domain.Users{},
		&domain.Drivers{},
		&domain.Vehicle{},
		&domain.DriverLocation{},
		&domain.Ride{},
		&domain.Payment{},
		&domain.Rating{},
		&domain.ChatMessage{},
		&domain.EmergencyAlert{},
	)
	if err != nil {
		log.Printf("migration warning/failed: %v\n", err)
		return err
	}

	log.Println("Database schema auto-migrated successfully (including chat, emergency alerts, and extensions)")
	return nil
}
