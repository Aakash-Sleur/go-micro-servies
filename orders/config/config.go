package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DB_URI         string
	PORT           string
	REDIS_URI      string
	REDIS_PASSWORD string
	JWT_SECRET     string
}

func Load() *Config {
	godotenv.Load(".env")
	return &Config{
		DB_URI:         getEnv("DB_URI", ""),
		PORT:           getEnv("PORT", "8083"),
		REDIS_URI:      getEnv("REDIS_URI", ""),
		REDIS_PASSWORD: getEnv("REDIS_PASSWORD", ""),
		JWT_SECRET:     getEnv("JWT_SECRET", ""),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return defaultValue
}
