package user

import (
	"auth/internal/model"
	"auth/internal/repository/user/converter"
	modelRepo "auth/internal/repository/user/model"
	"context"

	db "github.com/pavel-pj/platform_lib/pkg/client/db"

	sq "github.com/Masterminds/squirrel"
)

func (r *UserRepository) Get(ctx context.Context, id int64) (model.User, error) {

	sqlstr, args, err := psql.Select().
		From(tableName).
		Columns(idColumn, columnName, columnEmail, columnRole, columnCreatedAt, columnUpdatedAt).
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

	var user modelRepo.User
	err = r.db.DB().ScanOneContext(ctx, &user, q, args...)
	if err != nil {
		return model.User{}, err
	}

	return converter.ToServiceFromRepoUser(user), nil

}
