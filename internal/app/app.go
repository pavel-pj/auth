package app

import (
	env "auth/internal/config/env"
	desc "auth/pkg/user_v1"
	"context"
	"log"
	"net"
	"net/http"
	"sync"

	gwruntime "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pavel-pj/platform_lib/pkg/closer"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

// App хранит ВСЕ зависимости приложения
type App struct {
	//Services serv.UserService
	//Pool     *pgxpool.Pool
	serviceProvider *serviceProvider
	grpcServer      *grpc.Server
	httpServer      *http.Server
}

func NewApp(ctx context.Context) (*App, error) {
	// Создаем пул соединений pgx

	a := &App{}

	err := a.initDeps(ctx)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) Run() error {
	defer func() {
		closer.CloseAll()
		closer.Wait()
	}()

	var wg sync.WaitGroup

	wg.Add(2)
	go func() {
		defer wg.Done()
		err := a.runGRPCServer()
		if err != nil {
			log.Fatalf("failed to start grpc server: %v", err)
		}

	}()

	go func() {
		defer wg.Done()
		err := a.runHTTPServer()
		if err != nil {
			log.Fatalf("failed to start http server: %v", err)
		}

	}()

	wg.Wait()

	return nil

}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initConfig,
		a.initServiceProvider,
		a.initGRPCServer,
		a.initHTTPServer,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initConfig(_ context.Context) error {
	err := env.Load(".env")
	if err != nil {
		return err
	}

	return nil
}

func (a *App) initServiceProvider(_ context.Context) error {
	a.serviceProvider = NewServiceProvider()
	return nil
}

func (a *App) initGRPCServer(ctx context.Context) error {

	a.grpcServer = grpc.NewServer(
		grpc.Creds(insecure.NewCredentials()),
		grpc.ChainUnaryInterceptor(recoveryInterceptor),
	)
	// Только для тестовой среды - для удобства работы в Postman
	reflection.Register(a.grpcServer)

	desc.RegisterUserV1Server(a.grpcServer, a.serviceProvider.GetUserImplimentation(ctx))

	closer.Add(func() error {
		a.grpcServer.GracefulStop()
		return nil
	})

	return nil
}

func (a *App) initHTTPServer(ctx context.Context) error {

	mux := gwruntime.NewServeMux()

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	err := desc.RegisterUserV1HandlerFromEndpoint(ctx, mux, a.serviceProvider.GetGRPCConfig().Address(), opts)
	if err != nil {
		return err
	}

	a.httpServer = &http.Server{
		Addr:    a.serviceProvider.GetHTTPConfig().Address(),
		Handler: mux,
	}

	return nil

}

func (a *App) runGRPCServer() error {
	//log.Printf("GRPC Server is running in %s", a.serviceProvider.grpcConfig.Address())

	list, err := net.Listen("tcp", a.serviceProvider.GetGRPCConfig().Address())
	if err != nil {
		return err
	}
	log.Printf("GRPC Server is running in %s", a.serviceProvider.GetGRPCConfig().Address())
	err = a.grpcServer.Serve(list)
	if err != nil {
		return err
	}

	return nil

}

func (a *App) runHTTPServer() error {

	addr := a.serviceProvider.GetHTTPConfig().Address()
	log.Printf("HTTP Server is running in %s", addr)

	err := a.httpServer.ListenAndServe()

	if err != nil {
		return err
	}

	return nil
}
