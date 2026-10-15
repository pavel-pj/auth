package user

import (
	"auth/internal/model"
	converter "auth/internal/repository/user/converter"
	"context"
)

func (s UserService) Get(ctx context.Context, id int64) (model.User, error) {

	data, err := s.Repository.Get(ctx, id)
	if err != nil {
		return model.User{}, err
	}

	return converter.ToServiceFromRepoUser(data), nil

}
