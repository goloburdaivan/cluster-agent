package config

import (
	"crypto/rsa"
	"log"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

type Config struct {
	ApiURL       string
	JWTPublicKey *rsa.PublicKey
	RedisAddr    string
	RedisPass    string
	RedisDB      int
}

func NewConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on system envs")
	}

	pubKey := readPublicKey()

	return &Config{
		ApiURL:       os.Getenv("API_URL"),
		JWTPublicKey: pubKey,
		RedisAddr:    os.Getenv("REDIS_ADDR"),
		RedisPass:    os.Getenv("REDIS_PASS"),
		RedisDB:      0,
	}
}

func readPublicKey() *rsa.PublicKey {
	keyPath := os.Getenv("JWT_PUBLIC_KEY_PATH")
	if keyPath == "" {
		log.Fatal("JWT_PUBLIC_KEY_PATH is not set")
	}

	keyBytes, err := os.ReadFile(keyPath)
	if err != nil {
		log.Fatalf("Could not read public key file at %s: %v", keyPath, err)
	}

	pubKey, err := jwt.ParseRSAPublicKeyFromPEM(keyBytes)
	if err != nil {
		log.Fatalf("Invalid RSA public key: %v", err)
	}

	return pubKey
}
