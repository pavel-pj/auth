package log

import (
	"auth/internal/model"
	"auth/internal/repository/log/converter"
	"context"
)

func (s LogService) Store(ctx context.Context, log model.Log) error {
	data := converter.ToRepoLogFromService(log)
	return s.repository.Store(ctx, data)

}
