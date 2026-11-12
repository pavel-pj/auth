package auth

import (
	desc "auth/pkg/auth_v1"
	"context"
)

func (a *AuthImplementation) GetAccessToken(ctx context.Context, req *desc.GetAccessRequest) (*desc.GetAccessResponse, error) {

	refreshToken := req.GetRefreshToken()
	newAccessToken, err := a.serv.GetAccessToken(ctx, refreshToken)
	if err != nil {
		return &desc.GetAccessResponse{}, err
	}

	return &desc.GetAccessResponse{
		AccessToken: newAccessToken,
	}, nil
}
