package user

import (
	repo "github.com/pavel-pj/auth/internal/repository"
	serv "github.com/pavel-pj/auth/internal/service"

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
