package db

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type Handler func(ctx context.Context) error

// Интефес Client
type Client interface {
	DB() DB
	Close() error
}

// Интефес DB
type DB interface {
	SQLExecer
	Transactor
	Pinger
	Close()
}

type TxManager interface {
	ReadCommitted(ctx context.Context, f Handler) error
}

// Интефейс для работы с транзакциями
type Transactor interface {
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
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
	ExecContext(ctx context.Context, q Query, args ...any) (pgconn.CommandTag, error)
	QueryContext(ctx context.Context, q Query, args ...any) (pgx.Rows, error)
	QueryRowContext(ctx context.Context, q Query, args ...any) pgx.Row
}

type Pinger interface {
	Ping(ctx context.Context) error
}
