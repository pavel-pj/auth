package user

import (
	"context"
	converter "github.com/pavel-pj/auth/internal/converter/user"
	desc "github.com/pavel-pj/auth/pkg/user_v1"
)

func (s *UserImplementation) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {

	id := req.GetId()
	result, err := s.serv.Get(ctx, id)
	if err != nil {
		return &desc.GetResponse{}, err
	}

	return &desc.GetResponse{User: converter.ToProtoUserFromService(result)}, nil
}
