package app

import (
	"context"
	"flag"
	"fmt"
	grpcMiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/opentracing/opentracing-go"
	"github.com/semho/chat-microservices/chat-server/internal/closer"
	"github.com/semho/chat-microservices/chat-server/internal/config"
	"github.com/semho/chat-microservices/chat-server/internal/interceptor"
	"github.com/semho/chat-microservices/chat-server/internal/logger"
	"github.com/semho/chat-microservices/chat-server/internal/tracing"
	accessV1 "github.com/semho/chat-microservices/chat-server/pkg/access_v1"
	desc "github.com/semho/chat-microservices/chat-server/pkg/chat-server_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"
	"sync"
	"time"
)

type App struct {
	servicesProvider *serviceProvider
	grpcServer       *grpc.Server
	grpcConn         *grpc.ClientConn
	accessClient     accessV1.AccessV1Client
}

func NewApp(ctx context.Context) (*App, error) {
	a := &App{}
	err := a.initDeps(ctx)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) Run() error {
	defer func() {
		a.grpcConn.Close() // close grpc connection
		closer.CloseAll()
		closer.Wait()
	}()

	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()

		err := a.runGRPCServer()
		if err != nil {
			log.Fatalf("failed to start grpc server: %v", err)
		}
	}()

	go func() {
		//go tool pprof -http=:8081 http://localhost:6060/debug/pprof/heap - посмотреть память в ui
		log.Println("Starting pprof server on  http://localhost:6060/debug/pprof/")
		if err := http.ListenAndServe("localhost:6060", nil); err != nil {
			log.Fatalf("pprof server failed: %v", err)
		}
	}()

	wg.Wait()

	return nil
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initConfig,
		a.initServiceProvider,
		a.initJaegerConfig,
		a.InitLogger,
		a.initClientConfig,
		a.initGRPCClient,
		a.initGRPCServer,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

var configPath string
var logLevel string

func init() {
	flag.StringVar(&configPath, "config-path", ".env", "path to config file")
	flag.StringVar(&logLevel, "l", "info", "log level")
}

func (a *App) initConfig(_ context.Context) error {
	flag.Parse()
	err := config.Load(configPath)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) initJaegerConfig(_ context.Context) error {
	a.servicesProvider.JaegerConfig()
	return nil
}

func (a *App) InitLogger(_ context.Context) error {
	flag.Parse()
	err := logger.InitDefault(logLevel)
	if err != nil {
		return err
	}
	tracing.Init(logger.Logger(), "ChatService", a.servicesProvider.jaegerConfig.Address())

	return nil
}

func (a *App) initServiceProvider(_ context.Context) error {
	a.servicesProvider = newServiceProvider()
	return nil
}

func (a *App) initClientConfig(_ context.Context) error {
	a.servicesProvider.ClientConfig()
	return nil
}

func (a *App) initGRPCServer(ctx context.Context) error {
	a.grpcServer = grpc.NewServer(
		grpc.Creds(insecure.NewCredentials()),
		grpc.UnaryInterceptor(
			grpcMiddleware.ChainUnaryServer(
				interceptor.AuthInterceptor(a.accessClient),
				interceptor.ServerTracingInterceptor,
				interceptor.LogInterceptor,
			),
		),
	)

	reflection.Register(a.grpcServer)
	desc.RegisterChatServerV1Server(a.grpcServer, a.servicesProvider.GetChatServerImpl(ctx))

	return nil
}

func (a *App) runGRPCServer() error {
	if a.accessClient == nil {
		return fmt.Errorf("gRPC client is not initialized")
	}

	log.Printf("Starting gRPC server on port: %s", a.servicesProvider.GRPCConfig().Address())

	list, err := net.Listen("tcp", a.servicesProvider.GRPCConfig().Address())
	if err != nil {
		return err
	}

	err = a.grpcServer.Serve(list)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) initGRPCClient(_ context.Context) error {
	var err error

	backoffConfig := backoff.Config{
		BaseDelay:  1 * time.Second,
		Multiplier: 1.6,
		Jitter:     0.2,
		MaxDelay:   15 * time.Second,
	}

	a.grpcConn, err = grpc.Dial(
		a.servicesProvider.clientConfig.Address(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
		grpc.WithConnectParams(grpc.ConnectParams{Backoff: backoffConfig, MinConnectTimeout: 10 * time.Second}),
		grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer())),
	)
	if err != nil {
		return fmt.Errorf("failed to connect to server: %v", err)
	}

	a.accessClient = accessV1.NewAccessV1Client(a.grpcConn)
	return nil
}
