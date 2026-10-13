package user

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"

	repo "auth/internal/repository"
	desc "auth/pkg/user_v1"
)

type UserRepository struct {
	DB *pgxpool.Pool
}

var psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

func NewUserRepository(DB *pgxpool.Pool) repo.UserRepository {
	return &UserRepository{
		DB: DB,
	}
}

func (r *UserRepository) Store(ctx context.Context, users *desc.UserInfo) (int64, error) {
	return 25, nil
}
