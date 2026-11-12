package auth

import (
	"auth/internal/model"
	"auth/internal/service"
	"auth/internal/utils"
	"context"
	"errors"
)

func (s *AuthService) GetAccessToken(ctx context.Context, refreshToken string) (string, error) {

	claims, err := utils.VerifyToken(refreshToken, s.JWT.RefreshTokenSecret)
	if err != nil {
		return "", errors.New("Could not refresh token")
	}

	newAccessToken, err := utils.GenerateToken(model.UserAuthInfo{
		Name: claims.Username,
		Role: claims.Role,
	},
		s.JWT.AccessTokenSecret,
		service.AccessTokenExpiration,
	)
	if err != nil {
		return "", err
	}

	return newAccessToken, nil

}
