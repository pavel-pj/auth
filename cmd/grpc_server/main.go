package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/brianvoe/gofakeit"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/timestamppb"

	desc "chat/pkg/note_v1"
)

const grpcPort = 50301

type server struct {
	desc.UnimplementedNoteV1Server
}

func (s *server) Store(ctx context.Context, req *desc.StoreRequest) (*desc.StoreResponse, error) {

	d := req.GetNote()
	fmt.Println(d)

	return &desc.StoreResponse{Id: gofakeit.Int64()}, nil
}

func (s *server) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {

	return &desc.GetResponse{
		Note: &desc.Note{
			Id: req.GetId(),
			Info: &desc.NoteInfo{
				Title:    "Георгий Жуков",
				Content:  "a couple of sheet",
				Author:   "Prundel",
				IsPublic: true,
			},
			CreatedAt: timestamppb.New(gofakeit.Date()),
			UpdatedAt: timestamppb.New(gofakeit.Date()),
		},
	}, nil

}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatal("Server could not be run")
	}

	s := grpc.NewServer()
	reflection.Register(s)
	desc.RegisterNoteV1Server(s, &server{})

	log.Printf("Server listening at %d", grpcPort)
	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to server: ", err)
	}

}
