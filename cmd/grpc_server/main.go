package main

import (
	"auth/internal/app"
	"auth/internal/config/env"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	userApi "auth/internal/api/user"
	desc "auth/pkg/user_v1"
)

func main() {
	PGCfg, err := env.NewPGConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	GRPCCfg, err := env.NewGRPCConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Создаём приложение
	app, err := app.NewApp(PGCfg, GRPCCfg)
	if err != nil {
		log.Fatal("Failed to create app: ", err)
	}
	defer func() {
		if err := app.Close(); err != nil {
			log.Printf("Warning: failed to close app: %v", err)
		}
	}()

	log.Printf("✅ Successful connect to PostgreSQL")

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", GRPCCfg.Port()))
	if err != nil {
		log.Fatal("Server could not be run")
	}

	s := grpc.NewServer()
	reflection.Register(s)
	desc.RegisterUserV1Server(s, userApi.NewUserImplimentation(app.Services))

	log.Printf("✅ Server listening at %s", GRPCCfg.Port())
	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to server: ", err)
	}

}
