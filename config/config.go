package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppHost string
	AppPort string
}

func LoadConfig() *Config {
	godotenv.Load()

	return &Config{
		AppHost: getEnv("APP_HOST", "localhost"),
		AppPort: getEnv("APP_PORT", "8080"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
