package app

import (
	//"github.com/pavel-pj/auth/internal/closer"
	"context"
	"log"

	"github.com/pavel-pj/auth/internal/config"

	grpcAccessApi "github.com/pavel-pj/auth/internal/api/access"
	grpcAuthApi "github.com/pavel-pj/auth/internal/api/auth"
	grpcUserApi "github.com/pavel-pj/auth/internal/api/user"

	db "github.com/pavel-pj/platform_lib/pkg/client/db"
	pg "github.com/pavel-pj/platform_lib/pkg/client/db/pg"

	env "github.com/pavel-pj/auth/internal/config/env"

	repository "github.com/pavel-pj/auth/internal/repository"
	authRepo "github.com/pavel-pj/auth/internal/repository/auth"
	logRepo "github.com/pavel-pj/auth/internal/repository/log"
	userRepo "github.com/pavel-pj/auth/internal/repository/user"

	service "github.com/pavel-pj/auth/internal/service"
	accessServ "github.com/pavel-pj/auth/internal/service/access"
	authServ "github.com/pavel-pj/auth/internal/service/auth"
	logServ "github.com/pavel-pj/auth/internal/service/log"
	userServ "github.com/pavel-pj/auth/internal/service/user"

	"github.com/pavel-pj/platform_lib/pkg/client/db/transaction"
	"github.com/pavel-pj/platform_lib/pkg/closer"
)

type serviceProvider struct {
	pgConfig        config.PGConfig
	grpcConfig      config.GRPCConfig
	httpConfig      config.HTTPConfig
	swaggerConfig   config.SwaggerConfig
	vaultConfig     config.VaultConfig
	loggerConfig    config.LoggerConfig
	prometeusConfig config.PrometeusConfig

	//pgPool *pgxpool.Pool
	pgClient  db.Client
	txManager db.TxManager

	userRepository repository.UserRepository
	logRepository  repository.LogRepository
	authRepository repository.AuthRepository

	userService   service.UserService
	logService    service.LogService
	authService   service.AuthService
	accessService service.AccessService
	jwt           service.JWT

	userImpl   *grpcUserApi.UserImplementation
	authImpl   *grpcAuthApi.AuthImplementation
	accessImpl *grpcAccessApi.AccessImplementation
}

func NewServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (s *serviceProvider) SetJWTSecrets(access, refresh string) {
	s.jwt = service.JWT{
		AccessTokenSecret:  []byte(access),
		RefreshTokenSecret: []byte(refresh),
	}
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

func (s *serviceProvider) GetHTTPConfig() config.HTTPConfig {

	if s.httpConfig == nil {
		conf, err := env.NewHTTPConfig()
		if err != nil {
			log.Fatalf("failed to get http config: %s", err.Error())
		}

		s.httpConfig = conf
	}

	return s.httpConfig

}

func (s *serviceProvider) GetSwaggerConfig() config.SwaggerConfig {

	if s.swaggerConfig == nil {
		conf, err := env.NewSwaggerConfig()
		if err != nil {
			log.Fatalf("failed to get swagger config: %s", err.Error())
		}

		s.swaggerConfig = conf
	}

	return s.swaggerConfig

}

func (s *serviceProvider) GetVaultConfig() config.VaultConfig {
	if s.vaultConfig == nil {
		conf, err := env.NewVaultConfig()
		if err != nil {
			log.Fatalf("failed to get vault config: %s", err.Error())
		}
		s.vaultConfig = conf
	}

	return s.vaultConfig
}

func (s *serviceProvider) GetLoggerConfig() config.LoggerConfig {
	if s.loggerConfig == nil {
		conf, err := env.NewLoggerConfig()
		if err != nil {
			log.Fatalf("failed to get logger config: %s", err.Error())
		}
		s.loggerConfig = conf
	}

	return s.loggerConfig
}

func (s *serviceProvider) GetPrometeusConfig() config.PrometeusConfig {

	if s.prometeusConfig == nil {
		conf, err := env.NewPrometeusConfig()
		if err != nil {
			log.Fatalf("failed to get prometeus config: %s", err.Error())
		}

		s.prometeusConfig = conf
	}

	return s.prometeusConfig

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

func (s *serviceProvider) GetAuthRepository(ctx context.Context) repository.AuthRepository {
	if s.authRepository == nil {
		s.authRepository = authRepo.NewAuthRepository(s.GetDBClient(ctx))

	}
	return s.authRepository
}

func (s *serviceProvider) GetUserService(ctx context.Context) service.UserService {
	if s.userService == nil {
		s.userService = userServ.NewUserService(
			s.GetUserRepository(ctx),
			s.GetTxManager(ctx),
			s.GetLogService(ctx))
	}
	return s.userService
}

func (s *serviceProvider) GetLogService(ctx context.Context) service.LogService {
	if s.logService == nil {
		s.logService = logServ.NewLogService(s.GetLogRepository(ctx), s.GetTxManager(ctx))
	}

	return s.logService
}

func (s *serviceProvider) GetAuthService(ctx context.Context) service.AuthService {
	if s.authService == nil {
		s.authService = authServ.NewAuthService(
			s.GetAuthRepository(ctx),
			s.jwt,
		)
	}

	return s.authService
}

func (s *serviceProvider) GetAccessService(ctx context.Context) service.AccessService {
	if s.accessService == nil {
		s.accessService = accessServ.NewAccessService(
			s.jwt,
		)
	}

	return s.accessService
}

func (s *serviceProvider) GetUserImplimentation(ctx context.Context) *grpcUserApi.UserImplementation {
	if s.userImpl == nil {
		s.userImpl = grpcUserApi.NewUserImplimentation(s.GetUserService(ctx))
	}

	return s.userImpl
}

func (s *serviceProvider) GetAuthImplimentation(ctx context.Context) *grpcAuthApi.AuthImplementation {
	if s.authImpl == nil {
		s.authImpl = grpcAuthApi.NewAuthImplimentation(s.GetAuthService(ctx))
	}

	return s.authImpl
}

func (s *serviceProvider) GetAccessImplimentation(ctx context.Context) *grpcAccessApi.AccessImplementation {
	if s.accessImpl == nil {
		s.accessImpl = grpcAccessApi.NewAccessImplimentation(s.GetAccessService(ctx))
	}

	return s.accessImpl
}
