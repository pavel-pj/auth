package log

import (
	repo "auth/internal/repository"
	serv "auth/internal/service"

	db "github.com/pavel-pj/platform_lib/pkg/client/db"
)

type LogService struct {
	repository repo.LogRepository
	txManager  db.TxManager
}

func NewLogService(
	repository repo.LogRepository,
	txManager db.TxManager) serv.LogService {
	return LogService{
		repository: repository,
		txManager:  txManager,
	}
}
