package auth

import (
	"auth/internal/model"
	"auth/internal/service"
	"auth/internal/utils"
	"context"
)

func (s *AuthService) Login(ctx context.Context, login, password string) (string, error) {

	user, err := s.Repository.GetUserHash(ctx, login)
	if err != nil {
		return "", err
	}

	if ok := utils.VerifyPassword(password, user.Password); !ok {
		return "", errorWrongPassword
	}

	refreshToken, err := utils.GenerateToken(
		model.UserAuthInfo{
			Name: user.Name,
			Role: user.Role,
		},
		s.JWT.RefreshTokenSecret,
		service.RefreshTokenExpiration,
	)
	if err != nil {
		return "", err
	}

	//refreshToken, err := utils.GenerateToken(model.UserInfo{
	//		Username: login,
	//	})

	return refreshToken, nil
}
