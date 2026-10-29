package app

/*
import (
	apiUserInt "auth/internal/api/user"
	configInt "auth/internal/config"
	repoUserInt "auth/internal/repository/user"
	servUserInt "auth/internal/service/user"
	"log"

	config "auth/internal/config/env"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ServiceProvider struct {
	GrpcConfig configInt.GRPCConfig
	PgConfig   configInt.PGConfig

	PgPool *pgxpool.Pool

	ServiceUser servUserInt.UserService
	RepoUser    repoUserInt.UserRepository
	ApiUser     *apiUserInt.Implementation
}

func NewServiceProvider() *ServiceProvider {
	return &ServiceProvider{}
}

func (s *ServiceProvider) GetGRPCConfig() configInt.GRPCConfig {
	if s.GrpcConfig == nil {
		conf, err := config.NewGRPCConfig()
		if err != nil {
			log.Fatalf("Could not get config. Error: %s", err.Error())
		}

		s.GrpcConfig = conf

	}

	return s.GrpcConfig
}
*/

import (
	"auth/internal/closer"
	"auth/internal/config"
	"context"
	"log"

	grpcApi "auth/internal/api/user"
	implSerivce "auth/internal/api/user"
	db "auth/internal/client/db"
	pg "auth/internal/client/db/pg"
	"auth/internal/client/db/transaction"
	_ "auth/internal/closer"
	env "auth/internal/config/env"
	repository "auth/internal/repository"
	logRepo "auth/internal/repository/log"
	userRepo "auth/internal/repository/user"
	service "auth/internal/service"
	logServ "auth/internal/service/log"
	userServ "auth/internal/service/user"
)

type serviceProvider struct {
	pgConfig   config.PGConfig
	grpcConfig config.GRPCConfig

	//pgPool *pgxpool.Pool
	pgClient  db.Client
	txManager db.TxManager

	userRepository repository.UserRepository
	logRepository  repository.LogRepository

	userService service.UserService
	logService  service.LogService

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

func (s *serviceProvider) GetDBClient(ctx context.Context) db.Client {
	if s.pgClient == nil {

		client, err := pg.New(context.Background(), s.GetPGConfig().DSN())

		if err != nil {
			log.Fatalf("failed to connect to database: %s", err.Error())
		}

		// Проверяем подключение
		err = client.DB().Ping(context.Background())
		if err != nil {
			log.Fatalf("failed to ping database: %s", err.Error())
		}

		closer.Add(func() error {
			// тк pool.Close не возвращает error - создаем обертку
			client.Close()
			return nil
		})

		s.pgClient = client

	}

	return s.pgClient
}

func (s *serviceProvider) GetTxManager(ctx context.Context) db.TxManager {
	if s.txManager == nil {
		s.txManager = transaction.NewTransationManager(s.GetDBClient(ctx).DB())
	}

	return s.txManager
}

func (s *serviceProvider) GetUserRepository(ctx context.Context) repository.UserRepository {

	if s.userRepository == nil {
		s.userRepository = userRepo.NewUserRepository(s.GetDBClient(ctx))
	}

	return s.userRepository
}

func (s *serviceProvider) GetLogRepository(ctx context.Context) repository.LogRepository {
	if s.logRepository == nil {
		s.logRepository = logRepo.NewLogRepository(s.GetDBClient(ctx))

	}
	return s.logRepository
}

func (s *serviceProvider) GetUserService(ctx context.Context) service.UserService {
	if s.userService == nil {
		s.userService = userServ.NewUserService(s.GetUserRepository(ctx), s.GetTxManager(ctx), s.GetLogService(ctx))
	}
	return s.userService
}

func (s *serviceProvider) GetLogService(ctx context.Context) service.LogService {
	if s.logService == nil {
		s.logService = logServ.NewLogService(s.GetLogRepository(ctx), s.GetTxManager(ctx))
	}

	return s.logService
}

func (s *serviceProvider) GetUserImplimentation(ctx context.Context) *grpcApi.Implementation {
	if s.userImpl == nil {
		s.userImpl = implSerivce.NewUserImplimentation(s.GetUserService(ctx))
	}

	return s.userImpl
}
