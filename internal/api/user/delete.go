package user

import (
	"context"
	desc "github.com/pavel-pj/auth/pkg/user_v1"
)

func (s *UserImplementation) Delete(ctx context.Context, req *desc.DeleteRequest) (*desc.DeleteResponse, error) {

	id := req.GetId()
	err := s.serv.Delete(ctx, id)
	if err != nil {
		return &desc.DeleteResponse{}, err
	}

	return &desc.DeleteResponse{}, nil
}
