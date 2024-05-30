package tests

import (
	"context"
	"fmt"
	"github.com/brianvoe/gofakeit"
	"github.com/gojuno/minimock/v3"
	chatServerAPI "github.com/semho/chat-microservices/chat-server/internal/api/chat-server"
	"github.com/semho/chat-microservices/chat-server/internal/model"
	"github.com/semho/chat-microservices/chat-server/internal/service"
	serviceMocks "github.com/semho/chat-microservices/chat-server/internal/service/mocks"
	desc "github.com/semho/chat-microservices/chat-server/pkg/chat-server_v1"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"testing"
)

func TestImplementation_CreateChat(t *testing.T) {
	t.Parallel()
	type chatServerServiceMockFunc func(mc *minimock.Controller) service.ChatServerService
	type args struct {
		ctx context.Context
		req *desc.CreateChatRequest
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id    = gofakeit.Int64()
		name1 = gofakeit.Name()
		name2 = gofakeit.Name()

		req = &desc.CreateChatRequest{
			Usernames: []*desc.User{
				{Name: name1},
				{Name: name2},
			},
		}

		users = []*model.User{
			{Name: name1},
			{Name: name2},
		}

		res = &desc.CreateChatResponse{
			Id: id,
		}

		serviceErr = fmt.Errorf("service error")
	)

	tests := []struct {
		name string

		args                  args
		want                  *desc.CreateChatResponse
		err                   error
		chatServerServiceMock chatServerServiceMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: res,
			err:  nil,
			chatServerServiceMock: func(mc *minimock.Controller) service.ChatServerService {
				mock := serviceMocks.NewChatServerServiceMock(mc)
				mock.CreateChatMock.Expect(ctx, users).Return(id, nil)
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
				mock.CreateChatMock.Expect(ctx, users).Return(0, serviceErr)
				return mock
			},
		},
		{
			name: "error GetUsernames nil",
			args: args{
				ctx: ctx,
				req: nil,
			},
			want: nil,
			err:  status.Error(codes.InvalidArgument, "Invalid request: UserNames must be provided"),
			chatServerServiceMock: func(mc *minimock.Controller) service.ChatServerService {
				mock := serviceMocks.NewChatServerServiceMock(mc)
				return mock
			},
		},
	}
	initLogger()
	for _, tt := range tests {
		tt := tt
		t.Run(
			tt.name, func(t *testing.T) {
				t.Parallel()
				chatServerServiceMock := tt.chatServerServiceMock(mc)
				api := chatServerAPI.NewImplementation(chatServerServiceMock)

				resHandler, err := api.CreateChat(tt.args.ctx, tt.args.req)
				require.Equal(t, tt.err, err)
				require.Equal(t, tt.want, resHandler)
			},
		)
	}
}
