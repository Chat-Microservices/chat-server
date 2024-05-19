package chatServerAPI

import (
	"context"
	"github.com/semho/chat-microservices/chat-server/internal/converter"
	"github.com/semho/chat-microservices/chat-server/internal/logger"
	desc "github.com/semho/chat-microservices/chat-server/pkg/chat-server_v1"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (i *Implementation) CreateChat(ctx context.Context, req *desc.CreateChatRequest) (
	*desc.CreateChatResponse,
	error,
) {
	if req.GetUsernames() == nil {
		return nil, status.Error(codes.InvalidArgument, "Invalid request: UserNames must be provided")
	}
	logger.Info("Creating chat with names: ", zap.Any("usernames", req.GetUsernames()))

	chatID, err := i.chatServerService.CreateChat(ctx, converter.ToUserModelFromUserApi(req.GetUsernames()))
	if err != nil {
		return nil, err
	}
	newChat := &desc.Chat{
		Id:        chatID,
		Usernames: req.GetUsernames(),
	}
	logger.Info("New chat: ", zap.Any("chat", newChat))

	return &desc.CreateChatResponse{
		Id: newChat.Id,
	}, nil
}
