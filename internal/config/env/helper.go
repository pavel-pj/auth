package env

import (
	"errors"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func Load(path string) error {

	err := godotenv.Load(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			log.Printf("warning: %s not found, using environment variables", path)
			return nil
		}
		return err
	}

	return nil

}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
