package auth

import (
	repo "auth/internal/repository"

	db "github.com/pavel-pj/platform_lib/pkg/client/db"
)

type AuthRepository struct {
	db db.Client
}

func NewAuthRepository(db db.Client) repo.AuthRepository {
	return &AuthRepository{
		db: db,
	}
}
