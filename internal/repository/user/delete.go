package user

import (
	db "auth/internal/client/db"
	"context"

	sq "github.com/Masterminds/squirrel"
)

func (r *UserRepository) Delete(ctx context.Context, id int64) (int64, error) {

	sqlStr, args, err := psql.Delete(tableName).
		Where(sq.Eq{idColumn: id}).
		ToSql()

	if err != nil {
		return 0, err
	}

	q := db.Query{
		Name:     "user_repository.Delete",
		QueryRow: sqlStr,
	}

	rows, err := r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return 0, err
	}

	return rows.RowsAffected(), err
}
