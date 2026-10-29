package user

import (
	desc "auth/pkg/user_v1"
	"context"
)

func (s *Implementation) Delete(ctx context.Context, req *desc.DeleteRequest) (*desc.DeleteResponse, error) {

	id := req.GetId()
	err := s.serv.Delete(ctx, id)
	if err != nil {
		return &desc.DeleteResponse{}, err
	}

	return &desc.DeleteResponse{}, nil
}
