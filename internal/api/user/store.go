package user

import (
	"context"
	converter "github.com/pavel-pj/auth/internal/converter/user"
	desc "github.com/pavel-pj/auth/pkg/user_v1"
)

func (s *UserImplementation) Store(ctx context.Context, req *desc.StoreRequest) (*desc.StoreResponse, error) {

	if req.GetUser().Password != req.GetUser().PasswordConfirm {
		return &desc.StoreResponse{}, ErrorPasswwordConfirmation
	}

	user := converter.ToUserCreateFromProto(req.GetUser())

	result, err := s.serv.Store(ctx, user)
	if err != nil {
		return &desc.StoreResponse{}, err
	}

	return &desc.StoreResponse{Id: result}, nil
}
