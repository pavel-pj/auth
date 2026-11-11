package auth

import (
	"auth/internal/model"
	"auth/internal/utils"
	"context"
	"errors"
	"time"
)

func (s *AuthService) GetRefreshToken(ctx context.Context, refreshToken string) (string, error) {

	claims, err := utils.VerifyToken(refreshToken, s.JWTSecret)
	if err != nil {
		return "", errors.New("Could not refresh token")
	}

	newRefreshToken, err := utils.GenerateToken(model.UserAuthInfo{
		Name: claims.Username,
		Role: claims.Role,
	},
		s.JWTSecret,
		2*time.Minute,
	)
	if err != nil {
		return "", err
	}

	return newRefreshToken, nil

}
