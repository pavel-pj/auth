package user

import (
	converter "auth/internal/converter/user"
	desc "auth/pkg/user_v1"
	"context"
)

func (s *Implementation) Index(ctx context.Context, req *desc.IndexRequest) (*desc.IndexResponse, error) {

	limit := req.GetLimit()
	offset := req.GetOffset()
	data, err := s.serv.Index(ctx, limit, offset)
	if err != nil {
		return &desc.IndexResponse{}, err
	}

	var users []*desc.User

	for _, u := range data {
		users = append(users, converter.ToProtoUserFromService(u))
	}

	return &desc.IndexResponse{User: users}, nil
}
