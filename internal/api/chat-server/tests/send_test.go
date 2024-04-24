package tests

import (
	"context"
	"fmt"
	"github.com/brianvoe/gofakeit"
	"github.com/gojuno/minimock/v3"
	"github.com/semho/chat-microservices/chat-server/internal/api/chat-server"
	"github.com/semho/chat-microservices/chat-server/internal/model"
	"github.com/semho/chat-microservices/chat-server/internal/service"
	serviceMocks "github.com/semho/chat-microservices/chat-server/internal/service/mocks"
	desc "github.com/semho/chat-microservices/chat-server/pkg/chat-server_v1"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"testing"
	"time"
)

func TestImplementation_SendMessage(t *testing.T) {
	t.Parallel()
	type chatServerServiceMockFunc func(mc *minimock.Controller) service.ChatServerService
	type args struct {
		ctx context.Context
		req *desc.SendMessageRequest
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		username = gofakeit.Username()
		chatID   = gofakeit.Int64()
		text     = gofakeit.BeerName()
		created  = time.Now().UTC()

		message = &model.Message{
			UserName:  username,
			ChatID:    chatID,
			Text:      text,
			Timestamp: created,
		}

		req = &desc.SendMessageRequest{
			Id:        chatID,
			From:      username,
			Text:      text,
			Timestamp: timestamppb.New(created),
		}

		serviceErr = fmt.Errorf("service error")
	)

	tests := []struct {
		name                  string
		args                  args
		want                  *emptypb.Empty
		err                   error
		chatServerServiceMock chatServerServiceMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: &emptypb.Empty{},
			err:  nil,
			chatServerServiceMock: func(mc *minimock.Controller) service.ChatServerService {
				mock := serviceMocks.NewChatServerServiceMock(mc)
				mock.SendMessageMock.Expect(ctx, message).Return(nil)
				return mock
			},
		},
		{
			name: "error case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: nil,
			err:  serviceErr,
			chatServerServiceMock: func(mc *minimock.Controller) service.ChatServerService {
				mock := serviceMocks.NewChatServerServiceMock(mc)
				mock.SendMessageMock.Expect(ctx, message).Return(serviceErr)
				return mock
			},
		},
		{
			name: "error without id",
			args: args{
				ctx: ctx,
				req: &desc.SendMessageRequest{
					Id:        0,
					From:      username,
					Text:      text,
					Timestamp: timestamppb.New(created),
				},
			},
			want: nil,
			err:  status.Error(codes.InvalidArgument, "Invalid request: Id must be provided"),
			chatServerServiceMock: func(mc *minimock.Controller) service.ChatServerService {
				mock := serviceMocks.NewChatServerServiceMock(mc)
				return mock
			},
		},
		{
			name: "error empty text",
			args: args{
				ctx: ctx,
				req: &desc.SendMessageRequest{
					Id:        chatID,
					From:      username,
					Text:      "",
					Timestamp: timestamppb.New(created),
				},
			},
			want: nil,
			err:  status.Error(codes.InvalidArgument, "Invalid request: From or Text must be provided"),
			chatServerServiceMock: func(mc *minimock.Controller) service.ChatServerService {
				mock := serviceMocks.NewChatServerServiceMock(mc)
				return mock
			},
		},
		{
			name: "error empty from",
			args: args{
				ctx: ctx,
				req: &desc.SendMessageRequest{
					Id:        chatID,
					From:      "",
					Text:      text,
					Timestamp: timestamppb.New(created),
				},
			},
			want: nil,
			err:  status.Error(codes.InvalidArgument, "Invalid request: From or Text must be provided"),
			chatServerServiceMock: func(mc *minimock.Controller) service.ChatServerService {
				mock := serviceMocks.NewChatServerServiceMock(mc)
				return mock
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(
			tt.name, func(t *testing.T) {
				t.Parallel()
				chatServerServiceMock := tt.chatServerServiceMock(mc)
				api := chatServerAPI.NewImplementation(chatServerServiceMock)

				resHandler, err := api.SendMessage(tt.args.ctx, tt.args.req)
				require.Equal(t, tt.err, err)
				require.Equal(t, tt.want, resHandler)
			},
		)
	}
}
