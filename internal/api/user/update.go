package user

import (
	converter "auth/internal/converter/user"
	"auth/internal/model"
	desc "auth/pkg/user_v1"
	"context"
)

func (s *Implementation) Update(ctx context.Context, req *desc.UpdateRequest) (*desc.UpdateResponse, error) {

	//user := converter.ToUserPatchFromProto(req)
	return &desc.UpdateResponse{User: converter.ToProtoUserFromService(model.User{})}, nil

}
