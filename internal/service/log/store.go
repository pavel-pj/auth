package log

import (
	"context"
	"github.com/pavel-pj/auth/internal/model"
	"github.com/pavel-pj/auth/internal/repository/log/converter"
)

func (s LogService) Store(ctx context.Context, log model.Log) error {
	data := converter.ToRepoLogFromService(log)
	return s.repository.Store(ctx, data)

}
