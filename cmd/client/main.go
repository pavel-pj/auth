package main

import (
	descAccess "auth/pkg/access_v1"
	"context"
	"flag"
	"fmt"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

var accessToken = flag.String("a", "", "access token")

const servicePort = "localhost:50560"

func main() {

	flag.Parse()

	ctx := context.Background()
	md := metadata.New(map[string]string{"Authorization": "Bearer " + *accessToken})

	ctx = metadata.NewOutgoingContext(ctx, md)

	conn, err := grpc.NewClient(servicePort, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Could not create grpc client: %w", err)
	}

	defer conn.Close()

	cl := descAccess.NewAccessV1Client(conn)
	_, err = cl.Check(ctx, &descAccess.CheckRequest{
		EndpointAddress: "/user_v1.UserV1/Get",
	})

	if err != nil {
		log.Fatal(err.Error())

	}

	fmt.Println("Access granted")

}
