package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port             string
	DatabaseURL      string
	DatabasePort     string
	DatabaseUser     string
	DatabasePassword string
	DatabaseName     string
	APIVersion       string
	FinnhubAPI       string
	FinnhubWURL      string
}

func LoadConfig() *Config {
	_ = godotenv.Load()

	cfg := &Config{
		Port:             getEnv("PORT", "8080"),
		DatabaseURL:      getEnv("DATABASE_URL", "localhost"),
		DatabasePort:     getEnv("DB_PORT", "5432"),
		DatabaseUser:     getEnv("DB_USER", "postgre"),
		DatabasePassword: getEnv("DB_PASSWORD", ""),
		DatabaseName:     getEnv("DB_NAME", "postgre"),
		APIVersion:       getEnv("API_VERSION", "/api/v1"),
		FinnhubAPI:       getEnv("FINNHUB_API", ""),
		FinnhubWURL:      getEnv("FINNHUB_WURL", "wss://ws.finnhub.io"),
	}
	return cfg
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
