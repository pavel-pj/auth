package pg

import (
	"auth/internal/client/db"
	"context"
	"fmt"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type pg struct {
	dbc *pgxpool.Pool
}

func NewDB(dbc *pgxpool.Pool) db.DB {
	return &pg{
		dbc: dbc,
	}
}

func (p *pg) ScanOneContext(ctx context.Context, dest any, q db.Query, args ...any) error {
	//logQuery(ctx, q, args...)

	row, err := p.dbc.Query(ctx, q.QueryRow, args...)
	if err != nil {
		return err
	}

	return pgxscan.ScanOne(dest, row)
}

func (p *pg) ScanAllContext(ctx context.Context, dest any, q db.Query, args ...any) error {
	//logQuery(ctx, q, args...)

	rows, err := p.dbc.Query(ctx, q.QueryRow, args...)
	if err != nil {
		return err
	}

	return pgxscan.ScanAll(dest, rows)
}

func (p *pg) ExecContext(ctx context.Context, dest any, q db.Query, args ...any) (pgconn.CommandTag, error) {
	logQuery(ctx, q, args...)
	return p.dbc.Exec(ctx, q.QueryRow, args...)
}
func (p *pg) QueryContext(ctx context.Context, dest any, q db.Query, args ...any) (pgx.Rows, error) {
	logQuery(ctx, q, args...)
	return p.dbc.Query(ctx, q.QueryRow, args...)

}
func (p *pg) QueryRowContext(ctx context.Context, dest any, q db.Query, args ...any) pgx.Row {
	logQuery(ctx, q, args...)
	return p.dbc.QueryRow(ctx, q.QueryRow, args...)
}

func (p *pg) Ping(ctx context.Context) error {
	return p.dbc.Ping(ctx)
}

func (p *pg) Close() {
	p.dbc.Close()
}

func logQuery(ctx context.Context, q db.Query, args ...any) {
	//prettyQuery := prettier.Pretty(q.QueryRow, prettier)

	fmt.Println("SQL: %s", q.QueryRow)
}
