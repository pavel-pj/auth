package user

import (
	db "auth/internal/client/db"
	repo "auth/internal/repository"

	sq "github.com/Masterminds/squirrel"
)

type UserRepository struct {
	db db.Client
}

const (
	tableName       = "users"
	idColumn        = "id"
	columnName      = "name"
	columnEmail     = "email"
	columnPassword  = "password"
	columnRole      = "role"
	columnCreatedAt = "created_at"
	columnUpdatedAt = "updated_at"
)

var psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

func NewUserRepository(db db.Client) repo.UserRepository {
	return &UserRepository{
		db: db,
	}
}
