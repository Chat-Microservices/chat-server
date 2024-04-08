package chatServerAPI

import (
	"github.com/semho/chat-microservices/chat-server/internal/service"
	desc "github.com/semho/chat-microservices/chat-server/pkg/chat-server_v1"
)

type Implementation struct {
	desc.UnimplementedChatServerV1Server
	chatServerService service.ChatServerService
}

func NewImplementation(chatServerService service.ChatServerService) *Implementation {
	return &Implementation{
		chatServerService: chatServerService,
	}
}
