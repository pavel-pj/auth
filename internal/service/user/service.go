package user

import (
	repo "auth/internal/repository"
	serv "auth/internal/service"

	"github.com/pavel-pj/platform_lib/pkg/client/db"
)

type UserService struct {
	Repository repo.UserRepository
	txManager  db.TxManager
	logService serv.LogService
}

func NewUserService(
	repository repo.UserRepository,
	txManager db.TxManager,
	logService serv.LogService,
) serv.UserService {
	return UserService{
		Repository: repository,
		txManager:  txManager,
		logService: logService,
	}
}
