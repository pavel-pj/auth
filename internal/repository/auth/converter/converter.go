package converter

import (
	"github.com/pavel-pj/auth/internal/model"
	modelRepo "github.com/pavel-pj/auth/internal/repository/auth/model"
)

func UserLoginFromRepoToService(u modelRepo.UserLogin) model.UserLogin {
	return model.UserLogin{
		Name:     u.Name,
		Password: u.Password,
		Role:     u.Role,
	}
}
