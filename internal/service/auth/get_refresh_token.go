package auth

import (
	"auth/internal/model"
	"auth/internal/service"
	"auth/internal/utils"
	"context"
	"errors"
)

func (s *AuthService) GetRefreshToken(ctx context.Context, refreshToken string) (string, error) {

	claims, err := utils.VerifyToken(refreshToken, s.JWT.RefreshTokenSecret)
	if err != nil {
		return "", errors.New("Could not verify refresh token")
	}

	newRefreshToken, err := utils.GenerateToken(model.UserAuthInfo{
		Name: claims.Username,
		Role: claims.Role,
	},
		s.JWT.RefreshTokenSecret,
		service.RefreshTokenExpiration,
	)
	if err != nil {
		return "", err
	}

	return newRefreshToken, nil

}
