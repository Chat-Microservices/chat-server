package chatServerAPI

import (
	"context"
	"github.com/semho/chat-microservices/chat-server/internal/converter"
	desc "github.com/semho/chat-microservices/chat-server/pkg/chat-server_v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
)

func (i *Implementation) CreateChat(ctx context.Context, req *desc.CreateChatRequest) (*desc.CreateChatResponse, error) {
	if req.GetUsernames() == nil {
		return nil, status.Error(codes.InvalidArgument, "Invalid request: UserNames must be provided")
	}
	log.Printf("User names: %+v", req.GetUsernames())

	chatID, err := i.chatServerService.CreateChat(ctx, converter.ToUserModelFromUserApi(req.GetUsernames()))
	if err != nil {
		return nil, err
	}
	newChat := &desc.Chat{
		Id:        chatID,
		Usernames: req.GetUsernames(),
	}
	log.Printf("New chat: %+v", newChat)

	return &desc.CreateChatResponse{
		Id: newChat.Id,
	}, nil
}
