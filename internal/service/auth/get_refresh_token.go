package auth

import (
	"context"
	"errors"
	"github.com/pavel-pj/auth/internal/model"
	"github.com/pavel-pj/auth/internal/service"
	"github.com/pavel-pj/auth/internal/utils"
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
