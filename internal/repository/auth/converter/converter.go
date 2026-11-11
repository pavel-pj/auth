package converter

import (
	"auth/internal/model"
	modelRepo "auth/internal/repository/auth/model"
)

func UserLoginFromRepoToService(u modelRepo.UserLogin) model.UserLogin {
	return model.UserLogin{
		Name:     u.Name,
		Password: u.Password,
		Role:     u.Role,
	}
}
