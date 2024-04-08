package app

import (
	"context"
	"flag"
	"github.com/semho/chat-microservices/chat-server/internal/closer"
	"github.com/semho/chat-microservices/chat-server/internal/config"
	desc "github.com/semho/chat-microservices/chat-server/pkg/chat-server_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

type App struct {
	servicesProvider *serviceProvider
	grpcServer       *grpc.Server
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
			return err
		}
	}

	return nil
}

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", ".env", "path to config file")
}

func (a *App) initConfig(_ context.Context) error {
	flag.Parse()
	err := config.Load(configPath)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) initServiceProvider(_ context.Context) error {
	a.servicesProvider = newServiceProvider()
	return nil
}

func (a *App) initGRPCServer(ctx context.Context) error {
	a.grpcServer = grpc.NewServer(grpc.Creds(insecure.NewCredentials()))

	reflection.Register(a.grpcServer)
	desc.RegisterChatServerV1Server(a.grpcServer, a.servicesProvider.GetChatServerImpl(ctx))

	return nil
}

func (a *App) runGRPCServer() error {
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
