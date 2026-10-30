package user

import (
	"auth/internal/model"
	"fmt"

	"context"
	"errors"

	serv "auth/internal/service"

	"github.com/jackc/pgx/v5"
)

func (s UserService) Update(ctx context.Context, user model.UserPatch) (model.User, error) {

	var updatedUser model.User

	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {

		var errTx error
		updatedUser, errTx = s.Repository.Update(ctx, user)
		if errTx != nil {
			if errors.Is(errTx, pgx.ErrNoRows) {
				return serv.ErrNotFound
			}
			return nil
		}

		s.logService.Store(ctx, model.Log{
			Level:   "info",
			Message: fmt.Sprintf("user_update.%s", createUpdateLog(user)),
		})

		return nil
	})

	if err != nil {
		return model.User{}, err
	}

	return updatedUser, nil

}

func createUpdateLog(user model.UserPatch) string {

	var str string
	if user.Name != nil {
		str = "name: " + *user.Name + ";"
	}
	if user.Email != nil {
		str = str + "enail: " + *user.Email
	}
	if user.Roles != nil {
		str = str + "roles: " + *user.Roles
	}
	return str
}
