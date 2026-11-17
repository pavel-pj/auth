package main

import (
	"context"
	a "github.com/pavel-pj/auth/internal/app"
	"log"
)

// старт приложения
func main() {

	ctx := context.Background()

	app, err := a.NewApp(ctx)
	if err != nil {
		log.Fatalf("Fail to init app:%v", err)
	}

	err = app.Run()
	if err != nil {
		log.Fatalf("Fail to start app:%v", err)
	}

}
