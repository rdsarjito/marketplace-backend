package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	SecretKey string
	AppHost   string
	AppPort   string
}

func LoadConfig() *Config {
	return &Config{
		SecretKey: os.Getenv("SECRET_KEY"),
		AppHost:   os.Getenv("APP_HOST"),
		AppPort:   os.Getenv("APP_PORT"),
	}
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("Error loading .env file", err)
	}
}
