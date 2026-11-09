package main

import (
	desc "auth/pkg/user_v1"
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var (
	address = "localhost:50560"
	userID  = int64(12)
)

func main() {

	creds, err := credentials.NewClientTLSFromFile("ca.cert", "")
	if err != nil {
		log.Fatalf("Could not load tls sertificate: %w", err)
	}

	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatalf("Could not create grpc connection: %w", err)
	}

	defer conn.Close()

	c := desc.NewUserV1Client(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.Get(ctx, &desc.GetRequest{Id: userID})
	if err != nil {
		log.Fatalf("error get request: %w", err)
	}

	log.Printf("User info:", r.GetUser())

}
