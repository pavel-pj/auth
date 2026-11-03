package model

import "time"

type Level string

const (
	LevelDebug Level = "debug"
	LevelInfo  Level = "info"
	LevelWarn  Level = "warn"
	LevelError Level = "error"
)

// ff
type Log struct {
	ID        int64
	CreatedAt time.Time
	Level     Level
	Message   string
}
