package chatServerAPI

import (
	"context"
	"github.com/semho/chat-microservices/chat-server/internal/converter"
	"github.com/semho/chat-microservices/chat-server/internal/logger"
	desc "github.com/semho/chat-microservices/chat-server/pkg/chat-server_v1"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (i *Implementation) SendMessage(ctx context.Context, req *desc.SendMessageRequest) (*emptypb.Empty, error) {
	if req.GetId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "Invalid request: Id must be provided")
	}

	if req.GetFrom() == "" || req.GetText() == "" {
		return nil, status.Error(codes.InvalidArgument, "Invalid request: From or Text must be provided")
	}

	err := i.chatServerService.SendMessage(ctx, converter.ToMessageModelFromMessageApi(req))
	if err != nil {
		return nil, err
	}

	newMessage := &desc.SendMessageRequest{
		From: req.GetFrom(),
		Text: req.GetText(),
	}

	logger.Info("New message: ", zap.Any("message", newMessage))

	return &emptypb.Empty{}, nil
}
