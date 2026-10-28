package user

import (
	"context"
	"fmt"

	model "auth/internal/model"
	converter "auth/internal/repository/user/converter"

	"github.com/sppps/pwdhash"
)

func (s UserService) Store(ctx context.Context, user model.User) (int64, error) {

	var id int64

	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {

		hash, errTx := pwdhash.Hash(user.Password)
		if errTx != nil {
			//log.Fatal(errTx)
			return errTx
		}

		hashedUser := model.User{
			Name:     user.Name,
			Email:    user.Email,
			Password: hash,
			Roles:    user.Roles,
		}

		id, errTx = s.Repository.Store(ctx, converter.ToRepoUserFromService(hashedUser))
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
