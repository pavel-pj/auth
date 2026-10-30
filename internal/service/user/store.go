package user

import (
	"context"
	"fmt"

	model "auth/internal/model"
)

func (s UserService) Store(ctx context.Context, user model.UserCreate) (int64, error) {

	var id int64

	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {

		var errTx error
		id, errTx = s.Repository.Store(ctx, user)
		if errTx != nil {
			return errTx
		}

		errTx = s.logService.Store(ctx, model.Log{
			Level:   model.LevelInfo,
			Message: fmt.Sprintf("user_store.name:%s,email:%s.role:%s", user.Name, user.Email, user.Roles),
		})

		if errTx != nil {
			return errTx
		}

		return nil

	})

	if err != nil {
		return 0, err
	}

	return id, nil

}
