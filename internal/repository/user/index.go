package user

import (
	db "auth/internal/client/db"
	"auth/internal/model"
	"auth/internal/repository/user/converter"
	modelRepo "auth/internal/repository/user/model"
	"context"
	"log"
)

func (r *UserRepository) Index(ctx context.Context, limit, offset uint64) ([]model.User, error) {

	sqlStr, args, err := psql.Select().
		From(tableName).
		Columns(idColumn, columnName, columnEmail, columnRole, columnCreatedAt, columnUpdatedAt).
		OrderBy(idColumn).
		Limit(limit).
		Offset(offset).
		ToSql()

	if err != nil {
		log.Println("Repository.Could not create user stirng:", err)
		return []model.User{}, err
	}

	q := db.Query{
		Name:     "user_repository.index",
		QueryRow: sqlStr,
	}

	var usersRepo []modelRepo.User
	err = r.db.DB().ScanAllContext(ctx, &usersRepo, q, args...)
	if err != nil {
		return []model.User{}, err
	}

	var users []model.User

	for _, u := range usersRepo {
		users = append(users, converter.ToServiceFromRepoUser(u))
	}

	return users, nil

}
