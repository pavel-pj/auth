package user

import (
	"context"
	converter "github.com/pavel-pj/auth/internal/converter/user"
	desc "github.com/pavel-pj/auth/pkg/user_v1"
)

func (s *UserImplementation) Update(ctx context.Context, req *desc.UpdateRequest) (*desc.UpdateResponse, error) {

	data := converter.ToUserPatchFromProto(req)
	updated, err := s.serv.Update(ctx, data)
	if err != nil {
		return &desc.UpdateResponse{}, err
	}

	return &desc.UpdateResponse{User: converter.ToProtoUserFromService(updated)}, nil

}
