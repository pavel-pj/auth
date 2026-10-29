package user

import (
	db "auth/internal/client/db"
	model "auth/internal/repository/user/model"
	"context"
	"log"
)

func (r *UserRepository) Store(ctx context.Context, user model.User) (int64, error) {

	sqlStr, arg, err := psql.Insert(tableName).
		Columns(columnName, columnEmail, columnPassword, columnRole).
		Values(user.Name, user.Email, user.Password, user.Roles).
		Suffix("RETURNING ID").
		ToSql()

	if err != nil {
		log.Println("Repository.Could not create user stirng:", err)
		return 0, err
	}

	q := db.Query{
		Name:     "user_repository.Store",
		QueryRow: sqlStr,
	}

	var id int64
	err = r.db.DB().QueryRowContext(ctx, q, arg...).Scan(&id)
	if err != nil {
		log.Println("Repository.Could not create user:", err)
		return 0, err
	}

	return id, nil

}
