package chatServerAPI

import (
	"context"
	"github.com/semho/chat-microservices/chat-server/internal/logger"
	desc "github.com/semho/chat-microservices/chat-server/pkg/chat-server_v1"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (i *Implementation) DeleteChat(ctx context.Context, req *desc.DeleteRequest) (*emptypb.Empty, error) {
	if req.GetId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "Invalid request: Id must be provided")
	}
	logger.Warn("Delete chat by id: ", zap.Int64("id", req.GetId()))
	err := i.chatServerService.DeleteChat(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
