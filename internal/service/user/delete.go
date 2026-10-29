package user

import (
	"auth/internal/model"
	serv "auth/internal/service"
	"context"
	"fmt"
)

func (s UserService) Delete(ctx context.Context, id int64) error {

	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		rows, err := s.Repository.Delete(ctx, id)
		if err != nil {
			return err
		}
		fmt.Println("rows:%d", rows)
		if rows == 0 {
			return serv.ErrNotFound
		}

		errLog := s.logService.Store(ctx, model.Log{
			Level:   model.LevelInfo,
			Message: fmt.Sprintf("user_delete.id:%d", id),
		})
		if errLog != nil {
			return errLog
		}

		return nil

	})
	if err != nil {
		return err
	}

	return nil
}
