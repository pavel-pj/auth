package user

import (
	db "auth/internal/client/db"
	model "auth/internal/repository/user/model"
	"context"

	sq "github.com/Masterminds/squirrel"
)

func (r *UserRepository) Get(ctx context.Context, id int64) (model.User, error) {

	sqlstr, args, err := psql.Select().
		From(tableName).
		Columns(idColumn, columnName, columnEmail, columnRole).
		Where(sq.Eq{idColumn: id}).
		Limit(1).
		ToSql()
	if err != nil {
		return model.User{}, err
	}

	q := db.Query{
		Name:     "user_repository.Get",
		QueryRow: sqlstr,
	}

	var user model.User
	err = r.db.DB().ScanOneContext(ctx, &user, q, args...)
	if err != nil {
		return model.User{}, err
	}

	return user, nil

}
