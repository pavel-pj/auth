package user

import (
	"auth/internal/model"
	"context"
)

const (
	defaultLimit = 10
	maxLimit     = 100
	maxOffset    = 10_000
)

func (s UserService) Index(ctx context.Context, limit, offset uint64) ([]model.User, error) {

	if limit == 0 {
		limit = defaultLimit
	}
	if limit > maxLimit {
		limit = maxLimit
	}

	users, err := s.Repository.Index(ctx, limit, offset)
	if err != nil {
		return []model.User{}, err
	}

	return users, nil

}
