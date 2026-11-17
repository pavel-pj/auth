package auth

import (
	"context"
	desc "github.com/pavel-pj/auth/pkg/auth_v1"
)

func (s *AuthImplementation) Login(ctx context.Context, req *desc.LoginRequest) (*desc.LoginResponse, error) {

	token, err := s.serv.Login(ctx, req.GetUsername(), req.GetPassword())
	if err != nil {
		return &desc.LoginResponse{}, err
	}

	return &desc.LoginResponse{
		RefreshToken: token,
	}, nil
}
