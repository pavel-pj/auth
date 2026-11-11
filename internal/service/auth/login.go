package auth

import (
	"auth/internal/model"
	"auth/internal/utils"
	"context"
	"time"
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
		s.JWTSecret,
		2*time.Minute,
	)
	if err != nil {
		return "", err
	}

	//refreshToken, err := utils.GenerateToken(model.UserInfo{
	//		Username: login,
	//	})

	return refreshToken, nil
}
