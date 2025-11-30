package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppHost              string
	AppPort              string
	MidtransServerKey    string
	MidtransClientKey    string
	MidtransIsProduction bool
	FrontendURL          string // Frontend URL for payment redirect
}

func LoadConfig() *Config {
	godotenv.Load()

	return &Config{
		AppHost:              getEnv("APP_HOST", "localhost"),
		AppPort:              getEnv("APP_PORT", "8080"),
		MidtransServerKey:    getEnv("MIDTRANS_SERVER_KEY", ""),
		MidtransClientKey:    getEnv("MIDTRANS_CLIENT_KEY", ""),
		MidtransIsProduction: getEnvBool("MIDTRANS_IS_PRODUCTION", false),
		FrontendURL:          getEnv("FRONTEND_URL", "http://localhost:5173"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.ParseBool(value); err == nil {
			return parsed
		}
	}
	return defaultValue
}
