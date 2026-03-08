package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port              string
	DataBaseURL       string
	SupabaseURL       string
	SupabaseAnonKey   string
	SupabaseSericeKey string
}

func Load() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println(".env not found, using system env")
	}

	cfg := &Config{
		Port:              os.Getenv("PORT"),
		DataBaseURL:       os.Getenv(""),
		SupabaseURL:       os.Getenv(""),
		SupabaseAnonKey:   os.Getenv(""),
		SupabaseSericeKey: os.Getenv(""),
	}
	return cfg
}
