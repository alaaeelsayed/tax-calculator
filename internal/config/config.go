package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port      string
	TaxAPIURL string
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found, using environment variables and defaults")
	}

	return &Config{
		Port:      getEnv("PORT", "8080"),
		TaxAPIURL: getEnv("TAX_API_BASE_URL", "http://localhost:5001"),
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
