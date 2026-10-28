package model

import "time"

type Level string

const (
	LevelDebug Level = "debug"
	LevelInfo  Level = "info"
	LevelWarn  Level = "warn"
	LevelError Level = "error"
)

type Log struct {
	ID        int64     `db:"id"`
	CreatedAt time.Time `db:"created_at"`
	Level     Level     `db:"level"`
	Message   string    `db:"message"`
}
