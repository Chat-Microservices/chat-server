package chatServerAPI

import (
	"context"
	"github.com/semho/chat-microservices/chat-server/internal/converter"
	desc "github.com/semho/chat-microservices/chat-server/pkg/chat-server_v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
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

	log.Printf("New message: %+v", newMessage)

	return &emptypb.Empty{}, nil
}
