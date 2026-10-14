package main

import (
	"auth/internal/app"
	"auth/internal/config"
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	converter "auth/internal/converter/user"
	service "auth/internal/service"
	desc "auth/pkg/user_v1"
)

type server struct {
	desc.UnimplementedUserV1Server
	serv service.UserService
}

func (s *server) Store(ctx context.Context, req *desc.StoreRequest) (*desc.StoreResponse, error) {

	user := converter.ToUserFromProto(req.GetUser())
	result, err := s.serv.Store(ctx, user)
	if err != nil {
		return &desc.StoreResponse{}, err
	}
	//fmt.Println(data)

	return &desc.StoreResponse{Id: result}, nil
}

/*
func (s *server) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {

	return &desc.GetResponse{
		id: 1,
	}, nil

}*/

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Создаём приложение
	app, err := app.NewApp(*cfg)
	if err != nil {
		log.Fatal("Failed to create app: ", err)
	}
	defer func() {
		if err := app.Close(); err != nil {
			log.Printf("Warning: failed to close app: %v", err)
		}
	}()

	log.Printf("✅ Successful connect to PostgreSQL")

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.ServerPort))
	if err != nil {
		log.Fatal("Server could not be run")
	}

	s := grpc.NewServer()
	reflection.Register(s)
	desc.RegisterUserV1Server(s, &server{serv: app.Services})

	log.Printf("✅ Server listening at %s", cfg.ServerPort)
	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to server: ", err)
	}

}
