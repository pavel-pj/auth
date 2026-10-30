package user

import (
	"auth/internal/client/db"
	"auth/internal/model"
	"auth/internal/repository/user/converter"
	modelRepo "auth/internal/repository/user/model"
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
)

func (r *UserRepository) Update(ctx context.Context, user model.UserPatch) (model.User, error) {
	udp := psql.Update(tableName).Where(sq.Eq{idColumn: user.ID})

	if user.Name != nil {
		udp = udp.Set(columnName, user.Name)
	}
	if user.Email != nil {
		udp = udp.Set(columnEmail, user.Email)
	}
	if user.Roles != nil {
		udp = udp.Set(columnRole, user.Roles)
	}
	udp = udp.Set(columnUpdatedAt, time.Now())

	strSql, args, err := udp.Suffix("RETURNING id,name,email,created_at,updated_at").ToSql()
	if err != nil {
		return model.User{}, err
	}

	q := db.Query{
		Name:     "user_repository.Update",
		QueryRow: strSql,
	}

	var row modelRepo.User

	err = r.db.DB().ScanOneContext(ctx, &row, q, args...)
	if err != nil {

		return model.User{}, err

	}

	return converter.ToServiceFromRepoUser(row), err

}
