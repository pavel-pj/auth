package auth

import (
	desc "auth/pkg/auth_v1"
	"context"
)

func (a *AuthImplementation) GetRefreshToken(ctx context.Context, req *desc.GetRefreshRequest) (*desc.GetRefreshResponse, error) {

	refreshToken := req.GetRefreshToken()
	newRefreshToken, err := a.serv.GetRefreshToken(ctx, refreshToken)
	if err != nil {
		return &desc.GetRefreshResponse{}, err
	}

	return &desc.GetRefreshResponse{
		RefreshToken: newRefreshToken,
	}, nil
}
