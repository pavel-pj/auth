package app

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	env "github.com/pavel-pj/auth/internal/config/env"
	"github.com/pavel-pj/auth/internal/interceptor"
	"github.com/pavel-pj/auth/internal/vault"
	descAccess "github.com/pavel-pj/auth/pkg/access_v1"
	descAuth "github.com/pavel-pj/auth/pkg/auth_v1"
	_ "github.com/pavel-pj/auth/pkg/statik"
	desc "github.com/pavel-pj/auth/pkg/user_v1"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/rakyll/statik/fs"
	"github.com/rs/cors"

	gwruntime "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pavel-pj/auth/internal/logger"
	"github.com/pavel-pj/auth/internal/metric"
	"github.com/pavel-pj/platform_lib/pkg/closer"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
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
	swaggerServer   *http.Server
	prometeusServer *http.Server
	vault           *vault.Secret
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

	wg.Add(4)
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

	go func() {
		defer wg.Done()

		err := a.runSwaggerServer()
		if err != nil {
			log.Fatalf("failed to start swagger server: %v", err)
		}
	}()

	go func() {
		defer wg.Done()

		err := a.runPrometuesServer()
		if err != nil {
			log.Fatalf("failed to start prometeus server: %v", err)
		}
	}()

	wg.Wait()

	return nil

}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initConfig,
		a.initServiceProvider,
		a.initVault,
		a.initGRPCServer,
		a.initHTTPServer,
		a.initSwaggerServer,
		a.initPrometeus,
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

func (a *App) initVault(ctx context.Context) error {
	cfg := a.serviceProvider.GetVaultConfig()

	vaultCfg := vault.Config{
		Address:               cfg.Address(),
		Token:                 cfg.Token(),
		RoleID:                cfg.RoleID(),
		SecretID:              cfg.SecretID(),
		Mount:                 cfg.KVMount(),
		JWTSecretPath:         cfg.JWTSecretPath(),
		JWTAccessSecretField:  cfg.JWTAccessSecretField(),
		JWTRefreshSecretField: cfg.JWTRefreshSecretField(),

		PKIMount:      cfg.PKIMount(),
		PKIRole:       cfg.PKIRole(),
		PKICommonName: cfg.PKICommonName(),
		PKITTL:        cfg.PKITTL(),
		RootCAFile:    cfg.RootCAFile(),
	}

	jwt, err := vault.LoadJWT(ctx, vaultCfg)
	if err != nil {
		return fmt.Errorf("init vault jwt: %w", err)
	}

	tlsConfig, err := vault.LoadPKI(ctx, vaultCfg)
	if err != nil {
		return fmt.Errorf("init vault pki: %w", err)
	}

	a.vault = &vault.Secret{
		TLS: tlsConfig,
		JWT: jwt,
	}

	a.serviceProvider.SetJWTSecrets(
		jwt.AccessTokenSecret,
		jwt.RefreshTokenSecret,
	)

	log.Printf("TLS from Vault PKI and JWT secrets loaded")

	return nil

}

func (a *App) initGRPCServer(ctx context.Context) error {

	//creds, err := credentials.NewServerTLSFromFile("service.pem", "service.key")
	//if err != nil {
	//	log.Fatalf("Could not load tls sertificate:%w", err)
	//}

	//logger.Init(logger.GetCore(logger.GetAtomicLevel("debug")))

	err := metric.Init(ctx, a.serviceProvider.GetPrometeusConfig())
	if err != nil {
		return err
	}

	logger.Init(a.serviceProvider.GetLoggerConfig())

	a.grpcServer = grpc.NewServer(
		grpc.Creds(credentials.NewTLS(a.vault.TLS.ServerTLS)),
		//grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.ChainUnaryInterceptor(
			interceptor.RecoveryInterceptor,
			interceptor.MetricInterceptor,
			interceptor.ValidateInterceptor,
			interceptor.LoggerInterceptor,
		),
	)
	// Только для тестовой среды - для удобства работы в Postman
	reflection.Register(a.grpcServer)

	desc.RegisterUserV1Server(a.grpcServer, a.serviceProvider.GetUserImplimentation(ctx))
	descAuth.RegisterAuthV1Server(a.grpcServer, a.serviceProvider.GetAuthImplimentation(ctx))
	descAccess.RegisterAccessV1Server(a.grpcServer, a.serviceProvider.GetAccessImplimentation(ctx))

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

	corsMiddleWare := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"PUT", "POST", "GET", "PATCH", "DELETE"},
		AllowedHeaders:   []string{"Content-Type", "Accept", "Content-Length", "Authorization"},
		AllowCredentials: true,
	})
	a.httpServer = &http.Server{
		Addr:      a.serviceProvider.GetHTTPConfig().Address(),
		Handler:   corsMiddleWare.Handler(mux),
		TLSConfig: a.vault.TLS.ServerTLS.Clone(),
	}

	closer.Add(func() error {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return a.httpServer.Shutdown(ctx)

	})

	return nil

}

func (a *App) initSwaggerServer(_ context.Context) error {

	statikFS, err := fs.New()
	if err != nil {
		return err
	}

	mux := http.NewServeMux()

	mux.Handle("/", http.StripPrefix("/", http.FileServer(statikFS)))
	mux.HandleFunc("/api.swagger.json", serveSwaggerFile("/api.swagger.json"))

	a.swaggerServer = &http.Server{
		Addr:    a.serviceProvider.GetSwaggerConfig().Address(),
		Handler: mux,
	}

	return nil

}

func (a *App) runGRPCServer() error {

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

	err := a.httpServer.ListenAndServeTLS("", "")
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	if err != nil {
		return err
	}

	return nil
}

func (a *App) runSwaggerServer() error {

	addr := a.serviceProvider.GetSwaggerConfig().Address()
	log.Printf("Swagger Server is running in %s", addr)

	err := a.swaggerServer.ListenAndServe()

	if err != nil {
		return err
	}

	return nil
}

func serveSwaggerFile(path string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		log.Printf("Serving swagger file: %s", path)
		statikFS, err := fs.New()

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Printf("Open swagger file: %s", path)

		file, err := statikFS.Open(path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer file.Close()

		log.Printf("Read swagger file: %s", path)

		content, err := io.ReadAll(file)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Printf("Write swagger file: %s", path)

		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(content)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Printf("Served swagger file: %s", path)

	}
}

func (a *App) initPrometeus(_ context.Context) error {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	a.prometeusServer = &http.Server{
		Addr:    a.serviceProvider.GetPrometeusConfig().Address(),
		Handler: mux,
	}

	return nil
}

func (a *App) runPrometuesServer() error {
	log.Printf("Prometeus server is running on %s", a.prometeusServer.Addr)

	closer.Add(func() error {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()
		return a.prometeusServer.Shutdown(ctx)

	})
	return a.prometeusServer.ListenAndServe()

}
