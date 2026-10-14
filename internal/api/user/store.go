package user

import (
	converter "auth/internal/converter/user"
	desc "auth/pkg/user_v1"
	"context"
)

func (s *Implementation) Store(ctx context.Context, req *desc.StoreRequest) (*desc.StoreResponse, error) {

	user := converter.ToUserFromProto(req.GetUser())
	result, err := s.serv.Store(ctx, user)
	if err != nil {
		return &desc.StoreResponse{}, err
	}

	return &desc.StoreResponse{Id: result}, nil
}
