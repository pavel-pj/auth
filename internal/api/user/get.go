package user

import (
	"context"

	converter "github.com/pavel-pj/auth/internal/converter/user"
	desc "github.com/pavel-pj/auth/pkg/user_v1"
)

func (s *UserImplementation) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {

	/*
		time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
		d := rand.Intn(3)
		if d == 1 {
			return &desc.GetResponse{}, errors.New("ОШИБКА")
		}
	*/
	id := req.GetId()
	result, err := s.serv.Get(ctx, id)
	if err != nil {
		return &desc.GetResponse{}, err
	}

	return &desc.GetResponse{User: converter.ToProtoUserFromService(result)}, nil
}
