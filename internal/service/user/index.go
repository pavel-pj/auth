package user

import (
	"auth/internal/model"
	converter "auth/internal/repository/user/converter"
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

	data, err := s.Repository.Index(ctx, limit, offset)
	if err != nil {
		return []model.User{}, err
	}

	var users []model.User

	for _, u := range data {
		users = append(users, converter.ToServiceFromRepoUser(u))
	}

	return users, nil

}
