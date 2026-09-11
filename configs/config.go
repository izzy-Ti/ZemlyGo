package configs

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Port            string
	DataBaseURL     string
	NeonAuthJWKSURL string
	NeonAuthIssuer  string
	JWTSecret       string
	Env             string
	MapsAPIKey      string
	MapsProvider    string

	// Legacy fields (deprecated)
	SupabaseURL        string
	SupabaseAnonKey    string
	SupabaseServiceKey string
	SupabaseSericeKey  string
}

func Load() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println(".env not found, using system env")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5432/zemlygo?sslmode=disable"
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "zemlygo-super-secret-production-jwt-key-2026"
	}

	env := os.Getenv("ENV")
	if env == "" {
		env = "development"
	}

	neonAuthJWKSURL := os.Getenv("NEON_AUTH_JWKS_URL")
	neonAuthIssuer := os.Getenv("NEON_AUTH_ISSUER")

	mapsAPIKey := os.Getenv("MAPS_API_KEY")
	mapsProvider := os.Getenv("MAPS_PROVIDER")
	if mapsProvider == "" {
		mapsProvider = "osrm" // Default free OpenStreetMap / OSRM engine
	}

	cfg := &Config{
		Port:            port,
		DataBaseURL:     dbURL,
		NeonAuthJWKSURL: neonAuthJWKSURL,
		NeonAuthIssuer:  neonAuthIssuer,
		JWTSecret:       jwtSecret,
		Env:             env,
		MapsAPIKey:      mapsAPIKey,
		MapsProvider:    mapsProvider,
	}
	return cfg
}

func (c *Config) IsDevelopment() bool {
	return c.Env == "development" || c.Env == "dev" || c.Env == ""
}

func (c *Config) IsNeon() bool {
	return strings.Contains(c.DataBaseURL, "neon.tech")
}
