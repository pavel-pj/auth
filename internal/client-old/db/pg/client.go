package pg

import (
	"auth/internal/client/db"

	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

type pgClient struct {
	masterDBC db.DB
}

func New(ctx context.Context, dns string) (db.Client, error) {

	dbc, err := pgxpool.New(ctx, dns)
	if err != nil {
		log.Fatalf("failed to connect to database: %s", err.Error())
	}

	return &pgClient{
		masterDBC: &pg{dbc: dbc},
	}, nil
}

func (c *pgClient) DB() db.DB {

	return c.masterDBC
}

func (c *pgClient) Close() error {
	if c.masterDBC != nil {
		c.masterDBC.Close()
	}

	return nil
}
