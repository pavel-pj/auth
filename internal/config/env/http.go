package env

import (
	"log"
	"net"

	"github.com/joho/godotenv"
)

type HTTPConfig struct {
	host string
	port string
}

func NewHTTPConfig() (HTTPConfig, error) {
	// Загружаем .env файл
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	cfg := HTTPConfig{}
	cfg.host = getEnv("SERVER_HTTP_HOST", "")
	cfg.port = getEnv("SERVER_HTTP_PORT", "50570")

	return cfg, nil
}

func (s HTTPConfig) Address() string {
	return net.JoinHostPort(s.host, s.port)
}

func (s HTTPConfig) Port() string {
	return s.port
}
