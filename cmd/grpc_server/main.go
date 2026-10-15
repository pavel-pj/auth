package main

import (
	a "auth/internal/app"
	"context"
	"log"
)

func main() {
	ctx := context.Background()

	app, err := a.NewApp(ctx)
	if err != nil {
		log.Fatalf("Fail to init app:%w", err)
	}

	err = app.Run()
	if err != nil {
		log.Fatalf("Fail to start app:%w", err)
	}

}
