package app

import (
	"auth/internal/config"

	grpcApi "auth/internal/api/user"
	repository "auth/internal/repository/user"
	service "auth/internal/service/user"

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

func (s *serviceProvider) PGConfig() {
	/*
		if s.pgConfig == nil {
			cfg, err := config.NewPPConfig()
			if err != nil {
				log.Fatalf("failed to get pg config: %s", err.Error())
			}

		}

		s.pgConfig := cfg
	*/
	//return config.PGConfig{}
}
