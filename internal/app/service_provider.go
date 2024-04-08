package app

import (
	"context"
	chatServerAPI "github.com/semho/chat-microservices/chat-server/internal/api/chat-server"
	"github.com/semho/chat-microservices/chat-server/internal/client/db"
	"github.com/semho/chat-microservices/chat-server/internal/client/db/pg"
	"github.com/semho/chat-microservices/chat-server/internal/client/db/transaction"
	"github.com/semho/chat-microservices/chat-server/internal/closer"
	"github.com/semho/chat-microservices/chat-server/internal/config"
	"github.com/semho/chat-microservices/chat-server/internal/config/env"
	"github.com/semho/chat-microservices/chat-server/internal/repository"
	chatServerRepository "github.com/semho/chat-microservices/chat-server/internal/repository/chat-server"
	"github.com/semho/chat-microservices/chat-server/internal/service"
	chatServerService "github.com/semho/chat-microservices/chat-server/internal/service/chat-server"
	"log"
)

type serviceProvider struct {
	pgConfig   config.PGConfig
	grpcConfig config.GRPCConfig

	dbClient             db.Client
	txManger             db.TxManager
	chatServerRepository repository.ChatServerRepository

	chatServerService service.ChatServerService

	chatServerImpl *chatServerAPI.Implementation
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (s *serviceProvider) GetPGConfig() config.PGConfig {
	if s.pgConfig == nil {
		cfg, err := env.NewPGConfig()
		if err != nil {
			log.Fatalf("failed to get pg config: %v", err)
		}

		s.pgConfig = cfg
	}

	return s.pgConfig
}

func (s *serviceProvider) GRPCConfig() config.GRPCConfig {
	if s.grpcConfig == nil {
		cfg, err := env.NewGRPCConfig()
		if err != nil {
			log.Fatalf("failed to get grpc config: %v", err)
		}

		s.grpcConfig = cfg
	}

	return s.grpcConfig
}

func (s *serviceProvider) GetDBClient(ctx context.Context) db.Client {
	if s.dbClient == nil {
		cl, err := pg.New(ctx, s.GetPGConfig().DSN())
		if err != nil {
			log.Fatalf("failed to get db client: %v", err)
		}

		err = cl.DB().Ping(ctx)
		if err != nil {
			log.Fatalf("failed to ping db: %v", err)
		}

		closer.Add(cl.Close)

		s.dbClient = cl
	}

	return s.dbClient
}

func (s *serviceProvider) GetTxManager(ctx context.Context) db.TxManager {
	if s.txManger == nil {
		s.txManger = transaction.NewTransactionManager(s.GetDBClient(ctx).DB())
	}

	return s.txManger
}

func (s *serviceProvider) GetChatServerRepository(ctx context.Context) repository.ChatServerRepository {
	if s.chatServerRepository == nil {
		s.chatServerRepository = chatServerRepository.NewRepository(s.GetDBClient(ctx))
	}

	return s.chatServerRepository
}

func (s *serviceProvider) GetChatServerService(ctx context.Context) service.ChatServerService {
	if s.chatServerService == nil {
		s.chatServerService = chatServerService.NewService(s.GetChatServerRepository(ctx), s.GetTxManager(ctx))
	}

	return s.chatServerService
}

func (s *serviceProvider) GetChatServerImpl(ctx context.Context) *chatServerAPI.Implementation {
	if s.chatServerImpl == nil {
		s.chatServerImpl = chatServerAPI.NewImplementation(s.GetChatServerService(ctx))
	}

	return s.chatServerImpl
}
