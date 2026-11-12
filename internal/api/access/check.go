package auth

import (
	desc "auth/pkg/access_v1"
	"context"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (a *AccessImplementation) Check(ctx context.Context, req *desc.CheckRequest) (*emptypb.Empty, error) {

	path := req.GetEndpointAddress()
	err := a.serv.Check(ctx, path)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil

}
