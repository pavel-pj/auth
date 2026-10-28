package user

import (
	db "auth/internal/client/db"
	model "auth/internal/repository/user/model"
	"context"
	"log"

	sq "github.com/Masterminds/squirrel"

	repo "auth/internal/repository"
)

type UserRepository struct {
	db db.Client
}

const (
	tableName      = "users"
	idColumn       = "id"
	columnName     = "name"
	columnEmail    = "email"
	columnPassword = "password"
	columnRole     = "role"
)

var psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

func NewUserRepository(db db.Client) repo.UserRepository {
	return &UserRepository{
		db: db,
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

func (r *UserRepository) Get(ctx context.Context, id int64) (model.User, error) {

	sqlstr, args, err := psql.Select().
		From(tableName).
		Columns(idColumn, columnName, columnEmail, columnRole).
		Where(sq.Eq{idColumn: id}).
		ToSql()
	if err != nil {
		return model.User{}, err
	}

	q := db.Query{
		Name:     "user_repository.Get",
		QueryRow: sqlstr,
	}

	var user model.User
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Roles,
	)
	if err != nil {
		return model.User{}, err
	}

	return user, nil

}
