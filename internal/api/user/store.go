package user

import (
	converter "auth/internal/converter/user"
	desc "auth/pkg/user_v1"
	"context"
	"errors"
)

func (s *Implementation) Store(ctx context.Context, req *desc.StoreRequest) (*desc.StoreResponse, error) {

	if req.GetUser().Password != req.GetUser().PasswordConfirm {
		return &desc.StoreResponse{}, errors.New("passwords and password confirmation are not equal")
	}

	user := converter.ToUserCreateFromProto(req.GetUser())

	result, err := s.serv.Store(ctx, user)
	if err != nil {
		return &desc.StoreResponse{}, err
	}

	return &desc.StoreResponse{Id: result}, nil
}
