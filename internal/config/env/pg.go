package env

import (
	"log"

	"github.com/joho/godotenv"
)

type PGConfig struct {
	// Только для основной БД
	DBUser     string
	DBPassword string
	DBName     string
	DBHost     string
	DBPort     string
}

func NewPGConfig() (PGConfig, error) {
	// Загружаем .env файл
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	cfg := PGConfig{}

	cfg.DBUser = getEnv("POSTGRES_USER", "user")
	cfg.DBPassword = getEnv("POSTGRES_PASSWORD", "password")
	cfg.DBName = getEnv("POSTGRES_DB", "app")
	cfg.DBHost = getEnv("DB_HOST", "localhost")
	cfg.DBPort = getEnv("DB_PORT", "5432")

	return cfg, nil
}

func (c PGConfig) DSN() string {

	return "host=" + c.DBHost + " port=" + c.DBPort + " user=" + c.DBUser + " password=" + c.DBPassword + " dbname=" + c.DBName + " sslmode=disable"
}
