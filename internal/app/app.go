package app

import (
	"context"
	//"chat-server/internal/db/generated"
	//"chat-server/internal/services"
	config "auth/internal/config"
	userRepo "auth/internal/repository/user"
	serv "auth/internal/service"
	userService "auth/internal/service/user"

	//epository "auth/internal/grpc_server/repository"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
)

// App хранит ВСЕ зависимости приложения
type App struct {
	Services serv.UserService
	Pool     *pgxpool.Pool
}

func NewApp(cfg config.Config) (*App, error) {
	// Создаем пул соединений pgx

	dbURL := "host=" + cfg.DBHost + " port=" + cfg.DBPort + " user=" + cfg.DBUser + " password=" + cfg.DBPassword + " dbname=" + cfg.DBName + " sslmode=disable"

	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Проверяем подключение
	if err := pool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("SUPER ERROR failed to ping database: %w", err)
	}

	repository := userRepo.NewUserRepository(pool)
	service := userService.NewUserService(repository)

	return &App{
		Services: service,
		Pool:     pool,
	}, nil
}

// Close закрывает соединение с БД
func (a *App) Close() error {
	a.Pool.Close()
	return nil
}
