package transaction

import (
	"auth/internal/client/db"
	"auth/internal/client/db/pg"
	"context"
	"errors"

	"fmt"

	"github.com/jackc/pgx/v5"
)

type manager struct {
	db db.Transactor
}

func NewTransationManager(db db.Transactor) db.TxManager {

	return &manager{
		db: db,
	}
}

func (m *manager) transaction(ctx context.Context, opts pgx.TxOptions, fn db.Handler) (err error) {

	// Если транзакция вложенная , пропускаем инициацию новой транзакции и выполняем обработчи
	tx, ok := ctx.Value(pg.TxKey).(pgx.Tx)
	if ok {
		return fn(ctx)
	}

	tx, err = m.db.BeginTx(ctx, opts)
	if err != nil {
		return err
	}

	ctx = pg.MakeContextTx(ctx, tx)

	defer func() {

		if r := recover(); r != nil {
			err = fmt.Errorf("panic recovered: %s", r)
		}

		if err != nil {
			if errRollback := tx.Rollback(ctx); errRollback != nil {
				err = errors.Join(err, fmt.Errorf("rollback failed: %w", errRollback))
			}
			return

		}

		if nil == err {
			err = tx.Commit(ctx)
			if err != nil {
				err = fmt.Errorf("tx commit failed: %w", err)
			}
		}

	}()

	if err = fn(ctx); err != nil {
		err = fmt.Errorf("failed executing code inside transaction: %w", err)
	}

	return err

}

func (m *manager) ReadCommitted(ctx context.Context, fn db.Handler) error {
	txOpts := pgx.TxOptions{IsoLevel: pgx.ReadCommitted}
	return m.transaction(ctx, txOpts, fn)
}
