package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	// Только для основной БД
	DBUser     string
	DBPassword string
	DBName     string
	DBHost     string
	DBPort     string
	ServerPort string
}

func Load() (*Config, error) {
	// Загружаем .env файл
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	cfg := &Config{}

	cfg.DBUser = getEnv("POSTGRES_USER", "user")
	cfg.DBPassword = getEnv("POSTGRES_PASSWORD", "password")
	cfg.DBName = getEnv("POSTGRES_DB", "app")
	cfg.DBHost = getEnv("DB_HOST", "localhost")
	cfg.DBPort = getEnv("PG_PORT", "5432")

	cfg.ServerPort = getEnv("SERVER_PORT", "55555")
	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
