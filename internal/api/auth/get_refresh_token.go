package auth

import (
	desc "auth/pkg/auth_v1"
	"context"
)

func (a *AuthImplementation) GetRefreshToken(ctx context.Context, req *desc.GetRefreshRequest) (*desc.GetRefreshResponse, error) {

	refreshToken := req.GetRefreshToken()
	//newRefreshToken := a.serv.GetRefreshToken(ctx,refreshToken,)

	return &desc.GetRefreshResponse{
		RefreshToken: refreshToken,
	}, nil
}
