package env

import (
	"log"
	"net"

	"github.com/joho/godotenv"
)

type SwaggerConfig struct {
	host string
	port string
}

func NewSwaggerConfig() (SwaggerConfig, error) {
	// Загружаем .env файл
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	cfg := SwaggerConfig{}
	cfg.host = getEnv("SERVER_SWAGGER_HOST", "localhost")
	cfg.port = getEnv("SERVER_SWAGGER_PORT", "50580")

	return cfg, nil
}

func (s SwaggerConfig) Address() string {
	return net.JoinHostPort(s.host, s.port)
}

func (s SwaggerConfig) Port() string {
	return s.port
}
