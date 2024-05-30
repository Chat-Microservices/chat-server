package chatServer

import (
	"github.com/semho/chat-microservices/chat-server/internal/client/db"
	"github.com/semho/chat-microservices/chat-server/internal/client/rpc"
	"github.com/semho/chat-microservices/chat-server/internal/repository"
	"github.com/semho/chat-microservices/chat-server/internal/service"
)

type serv struct {
	chatServerRepository repository.ChatServerRepository
	txManager            db.TxManager
	authServiceClient    rpc.AuthServiceClient
}

func NewService(
	chatServerRepository repository.ChatServerRepository,
	txManager db.TxManager,
	authServiceClient rpc.AuthServiceClient,
) service.ChatServerService {
	return &serv{
		chatServerRepository: chatServerRepository,
		txManager:            txManager,
		authServiceClient:    authServiceClient,
	}
}
