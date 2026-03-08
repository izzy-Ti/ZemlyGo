package database

import (
	"log"

	"github.com/izzy-Ti/ZemlyGo/internals/domain"
	"gorm.io/gorm"
)

func RunMigration(db *gorm.DB) {
	err := db.AutoMigrate(
		&domain.Payment{},
		&domain.Rating{},
		&domain.Ride{},
		&domain.Vehcile{},
	)
	if err != nil {
		log.Fatal("migration failed: ", err)
	}

	log.Println("database migrated")
}
