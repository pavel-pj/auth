package app

import (
	"context"
	"log"
	"net"

	//"chat-server/internal/db/generated"
	//"chat-server/internal/services"

	//epository "auth/internal/grpc_server/repository"
	"google.golang.org/grpc"

	env "auth/internal/config/env"

	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	"auth/internal/closer"
	desc "auth/pkg/user_v1"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// App хранит ВСЕ зависимости приложения
type App struct {
	//Services serv.UserService
	//Pool     *pgxpool.Pool
	serviceProvider *serviceProvider
	grpcServer      *grpc.Server
}

func NewApp(ctx context.Context) (*App, error) {
	// Создаем пул соединений pgx

	a := &App{}

	err := a.initDeps(ctx)
	if err != nil {
		return nil, nil
	}

	return a, nil
}

func (a *App) Run() error {
	defer func() {
		closer.CloseAll()
		closer.Wait()
	}()

	return a.runGRPCServer()
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initConfig,
		a.initServiceProvider,
		a.initGRPCServer,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return nil
		}
	}

	return nil
}

func (a *App) initConfig(_ context.Context) error {
	err := env.Load(".env")
	if err != nil {
		return nil
	}
	return nil
}

func (a *App) initServiceProvider(_ context.Context) error {
	a.serviceProvider = NewServiceProvider()
	return nil
}

func (a *App) initGRPCServer(ctx context.Context) error {
	a.grpcServer = grpc.NewServer(grpc.Creds(insecure.NewCredentials()))

	reflection.Register(a.grpcServer)
	desc.RegisterUserV1Server(a.grpcServer, a.serviceProvider.GetUserImplimentation(ctx))

	return nil
}

func (a *App) runGRPCServer() error {
	//log.Printf("GRPC Server is running in %s", a.serviceProvider.grpcConfig.Address())

	list, err := net.Listen("tcp", a.serviceProvider.GetGRPCConfig().Address())
	if err != nil {
		return err
	}
	log.Printf("GRPC Server is running in %s", a.serviceProvider.grpcConfig.Address())
	err = a.grpcServer.Serve(list)
	if err != nil {
		return err
	}

	return nil

}
