package log

import (
	db "auth/internal/client/db"
	repo "auth/internal/repository"
	serv "auth/internal/service"
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
