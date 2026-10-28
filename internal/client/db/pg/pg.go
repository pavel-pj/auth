package pg

import (
	"auth/internal/client/db"
	"auth/internal/client/db/prettier"
	"context"
	"fmt"
	"log"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type pg struct {
	dbc *pgxpool.Pool
}

func (p *pg) ScanOneContext(ctx context.Context, dest any, q db.Query, args ...any) error {
	logQuery(ctx, q, args...)

	row, err := p.dbc.Query(ctx, q.QueryRow, args...)
	if err != nil {
		return err
	}

	return pgxscan.ScanOne(dest, row)
}

func (p *pg) ScanAllContext(ctx context.Context, dest any, q db.Query, args ...any) error {
	logQuery(ctx, q, args...)

	rows, err := p.dbc.Query(ctx, q.QueryRow, args...)
	if err != nil {
		return err
	}

	return pgxscan.ScanAll(dest, rows)
}

func (p *pg) ExecContext(ctx context.Context, q db.Query, args ...any) (pgconn.CommandTag, error) {
	logQuery(ctx, q, args...)
	return p.dbc.Exec(ctx, q.QueryRow, args...)
}
func (p *pg) QueryContext(ctx context.Context, q db.Query, args ...any) (pgx.Rows, error) {
	logQuery(ctx, q, args...)
	return p.dbc.Query(ctx, q.QueryRow, args...)

}
func (p *pg) QueryRowContext(ctx context.Context, q db.Query, args ...any) pgx.Row {
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

	//log.Printf("SQL [%s]: %s args=%v", q.Name, q.QueryRow, args)
	prettyQuery := prettier.Pretty(q.QueryRow, prettier.PlaceholderDollar, args...)

	log.Println(
		//ctx,
		fmt.Sprintf("sql: %s", q.Name),
		fmt.Sprintf("query: %s", prettyQuery),
	)
}
