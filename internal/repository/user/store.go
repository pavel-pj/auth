package user

import (
	db "auth/internal/client/db"
	"auth/internal/model"

	//modelRepo "auth/internal/repository/model"
	"context"
	"log"

	"github.com/sppps/pwdhash"
)

func (r *UserRepository) Store(ctx context.Context, user model.UserCreate) (int64, error) {

	hash, err := pwdhash.Hash(user.Password)
	if err != nil {
		return 0, err
	}

	sqlStr, arg, err := psql.Insert(tableName).
		Columns(columnName, columnEmail, columnPassword, columnRole).
		Values(user.Name, user.Email, hash, user.Roles).
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
