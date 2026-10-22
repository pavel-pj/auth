package env

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func Load(path string) error {

	err := godotenv.Load(path)
	if err != nil {
		return err
	}

	for _, kv := range os.Environ() {
		key, value, ok := strings.Cut(kv, "=")
		if ok {
			fmt.Printf("%s=%s\n", key, value)
		}
	}

	return nil

}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
