package bootstrap

import (
	"log"
	"strings"
	"time"

	"github.com/izzy-Ti/ZemlyGo/configs"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewDatabase(cfg *configs.Config) *gorm.DB {
	dbLogLevel := logger.Silent
	if cfg.IsDevelopment() {
		dbLogLevel = logger.Warn
	}

	isNeon := strings.Contains(cfg.DataBaseURL, "neon.tech")
	if isNeon {
		log.Println("[Database] Initializing connection to Neon Serverless PostgreSQL...")
		if strings.Contains(cfg.DataBaseURL, "sslmode=disable") {
			log.Println("[Database Warning] Neon requires SSL. Please ensure your connection string has 'sslmode=require'")
		}
	}

	db, err := gorm.Open(postgres.Open(cfg.DataBaseURL), &gorm.Config{
		Logger: logger.Default.LogMode(dbLogLevel),
	})
	if err != nil {
		log.Printf("\n===============================================================\n")
		log.Printf("[DATABASE ERROR] Failed to connect to database at: %s\n", cfg.DataBaseURL)
		log.Printf("Error: %v\n", err)
		log.Printf("Tips for Neon Database:\n")
		log.Printf("1. Log in to Neon console (https://console.neon.tech) and copy your connection string.\n")
		log.Printf("2. Set DATABASE_URL in your .env file, for example:\n")
		log.Printf("   DATABASE_URL=\"postgresql://neondb_owner:password@ep-xyz-pooler.us-east-2.aws.neon.tech/neondb?sslmode=require\"\n")
		log.Printf("3. Make sure to use the pooled connection string (with '-pooler' in the host) for best performance.\n")
		log.Printf("===============================================================\n\n")
		log.Fatal("Database connection required to start ZemlyGo server: ", err)
	}

	// Optimize connection pool for serverless PostgreSQL / Neon
	sqlDB, err := db.DB()
	if err == nil {
		// Serverless Neon suspends idle compute; keep connection lifetime reasonable
		// so stale connections are safely recycled before proxy drops them.
		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetMaxOpenConns(50)
		sqlDB.SetConnMaxLifetime(10 * time.Minute)
		sqlDB.SetConnMaxIdleTime(2 * time.Minute)
	}

	if isNeon {
		log.Println("[Database] Connected successfully to Neon Serverless PostgreSQL with optimized connection pooling!")
	} else {
		log.Println("[Database] Connected successfully to PostgreSQL")
	}

	return db
}
