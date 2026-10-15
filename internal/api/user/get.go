package user

import (
	converter "auth/internal/converter/user"
	desc "auth/pkg/user_v1"
	"context"
)

func (s *Implementation) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {

	id := req.GetId()
	result, err := s.serv.Get(ctx, id)
	if err != nil {
		return &desc.GetResponse{}, err
	}

	return &desc.GetResponse{User: converter.ToProtoUserFromService(result)}, nil
}
