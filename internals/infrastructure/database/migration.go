package database

import (
	"log"

	"gorm.io/gorm"
)

func RunMigration(db *gorm.DB) {
	err := db.AutoMigrate()
	if err != nil {
		log.Fatal("migration failed: ", err)
	}

	log.Println("database migrated")
}
