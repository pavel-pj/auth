package user

import (
	"context"
	"log"

	model "auth/internal/repository/user/model"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"

	repo "auth/internal/repository"
)

type UserRepository struct {
	DB *pgxpool.Pool
}

const (
	tableName      = "users"
	columnName     = "name"
	columnEmail    = "email"
	columnPassword = "password"
	columnRole     = "role"
)

var psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

func NewUserRepository(DB *pgxpool.Pool) repo.UserRepository {
	return &UserRepository{
		DB: DB,
	}
}

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

	var id int64
	err = r.DB.QueryRow(ctx, sqlStr, arg...).Scan(&id)
	if err != nil {
		log.Println("Repository.Could not create user:", err)
		return 0, err
	}

	return id, nil

}
