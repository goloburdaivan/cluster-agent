package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	ApiURL string
}

func NewConfig() *Config {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return &Config{
		ApiURL: os.Getenv("API_URL"),
	}
}
