package db

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type Client interface {
	DB()
	Close() error
}

type Query struct {
	Name     string
	QueryRow string
}

type SQLExecer interface {
	NamedExecer
	QueryExecer
}

type NamedExecer interface {
	ScanOneContext(ctx context.Context, dest any, q Query, args ...any) error
	ScanAllContext(ctx context.Context, dest any, q Query, args ...any) error
}

type QueryExecer interface {
	ExecContext(ctx context.Context, dest any, q Query, args ...any) (pgconn.CommandTag, error)
	QueryContext(ctx context.Context, dest any, q Query, args ...any) (pgx.Rows, error)
	QueryRowContext(ctx context.Context, dest any, q Query, args ...any) pgx.Rows
}

type Pinger interface {
	Ping(ctx context.Context) error
}

type DB interface {
	SQLExecer
	Pinger
	Close()
}
