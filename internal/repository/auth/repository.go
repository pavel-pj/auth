package auth

import (
	repo "auth/internal/repository"

	sq "github.com/Masterminds/squirrel"
	db "github.com/pavel-pj/platform_lib/pkg/client/db"
)

type AuthRepository struct {
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

func NewAuthRepository(db db.Client) repo.AuthRepository {
	return &AuthRepository{
		db: db,
	}
}
