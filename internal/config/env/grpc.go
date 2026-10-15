package env

import (
	"log"
	"net"

	"github.com/joho/godotenv"
)

type GRPCConfig struct {
	host string
	port string
}

func NewGRPCConfig() (GRPCConfig, error) {
	// Загружаем .env файл
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	cfg := GRPCConfig{}
	cfg.port = getEnv("SERVER_HOST", "")
	cfg.port = getEnv("SERVER_PORT", "88888")

	return cfg, nil
}

func (s GRPCConfig) Address() string {
	return net.JoinHostPort(s.host, s.port)
}

func (s GRPCConfig) Port() string {
	return s.port
}
