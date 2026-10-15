package converter

import (
	modelService "auth/internal/model"
	modelRepo "auth/internal/repository/user/model"
)

func ToRepoUserFromService(data modelService.User) modelRepo.User {

	return modelRepo.User{
		Name:     data.Name,
		Email:    data.Email,
		Password: data.Password,
		Roles:    data.Roles,
	}

}

func ToServiceFromRepoUser(data modelRepo.User) modelService.User {
	return modelService.User{
		ID:    data.ID,
		Name:  data.Name,
		Email: data.Email,
		Roles: data.Roles,
	}
}
