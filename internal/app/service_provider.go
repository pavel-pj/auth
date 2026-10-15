package app

import (
	"auth/internal/closer"
	"auth/internal/config"
	"context"
	"log"

	grpcApi "auth/internal/api/user"
	implSerivce "auth/internal/api/user"
	_ "auth/internal/closer"
	env "auth/internal/config/env"
	repository "auth/internal/repository"
	userRepo "auth/internal/repository/user"
	service "auth/internal/service"
	userServ "auth/internal/service/user"

	"github.com/jackc/pgx/v5/pgxpool"
)

type serviceProvider struct {
	pgConfig   config.PGConfig
	grpcConfig config.GRPCConfig

	pgPool *pgxpool.Pool

	userReoisitory repository.UserRepository
	userService    service.UserService

	userImpl *grpcApi.Implementation
}

func NewServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (s *serviceProvider) GetPGConfig() config.PGConfig {

	if s.pgConfig == nil {
		conf, err := env.NewPGConfig()
		if err != nil {
			log.Fatalf("failed to get pg config: %s", err.Error())
		}

		s.pgConfig = conf
	}

	return s.pgConfig

}

func (s *serviceProvider) GetGRPCConfig() config.GRPCConfig {

	if s.grpcConfig == nil {
		conf, err := env.NewGRPCConfig()
		if err != nil {
			log.Fatalf("failed to get grpc config: %s", err.Error())
		}
		s.grpcConfig = conf
	}

	return s.grpcConfig
}

func (s *serviceProvider) GetPgPool(ctx context.Context) *pgxpool.Pool {
	if s.pgPool == nil {

		pool, err := pgxpool.New(context.Background(), s.GetPGConfig().DSN())
		if err != nil {
			log.Fatalf("failed to connect to database: %s", err.Error())
		}

		// Проверяем подключение
		err = pool.Ping(context.Background())
		if err != nil {
			log.Fatalf("failed to ping database: %s", err.Error())
		}

		closer.Add(func() error {
			// тк pool.Close не возвращает error - создаем обертку
			pool.Close()
			return nil
		})

		s.pgPool = pool

	}

	return s.pgPool
}

func (s *serviceProvider) GetUserReoisitory(ctx context.Context) repository.UserRepository {
	if s.userReoisitory == nil {
		s.userReoisitory = userRepo.NewUserRepository(s.GetPgPool(ctx))
	}
	return s.userReoisitory
}

func (s *serviceProvider) GetUserService(ctx context.Context) service.UserService {
	if s.userService == nil {
		s.userService = userServ.NewUserService(s.GetUserReoisitory(ctx))
	}
	return s.userService
}

func (s *serviceProvider) GetUserImplimentation(ctx context.Context) *grpcApi.Implementation {
	if s.userImpl == nil {
		s.userImpl = implSerivce.NewUserImplimentation(s.GetUserService(ctx))
	}

	return s.userImpl
}
