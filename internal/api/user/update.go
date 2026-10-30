package user

import (
	converter "auth/internal/converter/user"
	desc "auth/pkg/user_v1"
	"context"
)

func (s *Implementation) Update(ctx context.Context, req *desc.UpdateRequest) (*desc.UpdateResponse, error) {

	data := converter.ToUserPatchFromProto(req)
	updated, err := s.serv.Update(ctx, data)
	if err != nil {
		return &desc.UpdateResponse{}, err
	}

	return &desc.UpdateResponse{User: converter.ToProtoUserFromService(updated)}, nil

}
