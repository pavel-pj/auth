package log

import (
	repo "github.com/pavel-pj/auth/internal/repository"
	serv "github.com/pavel-pj/auth/internal/service"

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
