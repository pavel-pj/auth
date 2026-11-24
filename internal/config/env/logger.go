package env

import (
	"log"
	"strconv"

	"github.com/joho/godotenv"
)

type LoggerConfig struct {
	level      string
	filename   string
	maxSize    int
	maxBackups int
	maxAge     int
}

func NewLoggerConfig() (LoggerConfig, error) {

	// Загружаем .env файл
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	cfg := LoggerConfig{}

	cfg.level = getEnv("LOGGER_LEVEL", "info")
	cfg.filename = getEnv("LOGGER_FILE_NAME", "logs/app.log")

	maxSize, err := strconv.Atoi(getEnv("LOGGER_MAX_SIZE", "10"))
	if err != nil {
		return LoggerConfig{}, err
	}
	cfg.maxSize = maxSize

	maxBackups, err := strconv.Atoi(getEnv("LOGGER_MAX_BACKUPS", "5"))
	if err != nil {
		return LoggerConfig{}, err
	}
	cfg.maxBackups = maxBackups

	maxAge, err := strconv.Atoi(getEnv("LOGGER_MAX_AGE", "7"))
	if err != nil {
		return LoggerConfig{}, err
	}
	cfg.maxAge = maxAge

	return cfg, nil
}

func (l LoggerConfig) Level() string {
	return l.level
}

func (l LoggerConfig) Filename() string {
	return l.filename
}
func (l LoggerConfig) MaxSize() int {
	return l.maxSize
}
func (l LoggerConfig) MaxBackups() int {
	return l.maxBackups
}
func (l LoggerConfig) MaxAge() int {
	return l.maxAge
}
