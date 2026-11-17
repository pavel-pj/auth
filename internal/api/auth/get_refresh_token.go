package auth

import (
	"context"
	desc "github.com/pavel-pj/auth/pkg/auth_v1"
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
