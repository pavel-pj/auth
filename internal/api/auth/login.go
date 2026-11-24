package auth

import (
	"context"

	"github.com/pavel-pj/auth/internal/logger"
	desc "github.com/pavel-pj/auth/pkg/auth_v1"
	"go.uber.org/zap"
)

func (s *AuthImplementation) Login(ctx context.Context, req *desc.LoginRequest) (*desc.LoginResponse, error) {

	logger.Debug("Debug. ",
		zap.Any("path", "API.Login."),
		zap.Any("req", req),
	)

	token, err := s.serv.Login(ctx, req.GetUsername(), req.GetPassword())
	if err != nil {
		return &desc.LoginResponse{}, err
	}

	return &desc.LoginResponse{
		RefreshToken: token,
	}, nil
}
