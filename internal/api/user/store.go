package user

import (
	converter "auth/internal/converter/user"
	desc "auth/pkg/user_v1"
	"context"
)

func (s *Implementation) Store(ctx context.Context, req *desc.StoreRequest) (*desc.StoreResponse, error) {

	err := req.Validate()
	if err != nil {
		return nil, err
	}

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
