package env

import (
	"log"
	"net"

	"github.com/joho/godotenv"
)

type PrometeusConfig struct {
	host      string
	port      string
	namespace string
	appName   string
}

func NewPrometeusConfig() (PrometeusConfig, error) {
	// Загружаем .env файл
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	cfg := PrometeusConfig{}
	cfg.host = getEnv("PROMETEUS_HTTP_HOST", "")
	cfg.port = getEnv("PROMETEUS_HTTP_PORT", "6060")
	cfg.namespace = getEnv("PROMETEUS_NAMESPACE", "test_namespace")
	cfg.appName = getEnv("PROMETEUS_APP_NAME", "test_app_name")

	return cfg, nil
}

func (s PrometeusConfig) Address() string {
	return net.JoinHostPort(s.host, s.port)
}

func (s PrometeusConfig) Namespace() string {
	return s.namespace
}
func (s PrometeusConfig) AppName() string {
	return s.appName
}
