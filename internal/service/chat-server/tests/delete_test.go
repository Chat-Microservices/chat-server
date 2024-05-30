package tests

import (
	"context"
	"errors"
	"fmt"
	"github.com/brianvoe/gofakeit"
	"github.com/gojuno/minimock/v3"
	"github.com/semho/chat-microservices/chat-server/internal/client/db"
	"github.com/semho/chat-microservices/chat-server/internal/client/rpc"
	"github.com/semho/chat-microservices/chat-server/internal/model"
	"github.com/semho/chat-microservices/chat-server/internal/repository"
	repoMocks "github.com/semho/chat-microservices/chat-server/internal/repository/mocks"
	chatServerService "github.com/semho/chat-microservices/chat-server/internal/service/chat-server"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_serv_DeleteChat(t *testing.T) {
	t.Parallel()
	type chatServerRepositoryMockFunc func(mc *minimock.Controller) repository.ChatServerRepository

	txManagerMock := func() db.TxManager {
		return &mockTxManager{}
	}
	authServiceClientMock := func() rpc.AuthServiceClient {
		return &mockAuthServiceClient{}
	}
	type args struct {
		ctx    context.Context
		chatId int64
	}

	var (
		ctx    = context.Background()
		mc     = minimock.NewController(t)
		chatId = gofakeit.Int64()

		qName = "CreateLog"
		qRow  = "INSERT INTO..."
		query = db.Query{
			Name:     qName,
			QueryRow: qRow,
		}
		log = &model.Log{
			Action:   qName,
			EntityId: chatId,
			Query:    qRow,
		}
		repoErr    = fmt.Errorf("repo error")
		repoErrLog = fmt.Errorf("repo error log")
	)

	tests := []struct {
		name                     string
		args                     args
		err                      error
		chatServerRepositoryMock chatServerRepositoryMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx:    ctx,
				chatId: chatId,
			},
			err: nil,
			chatServerRepositoryMock: func(mc *minimock.Controller) repository.ChatServerRepository {
				mock := repoMocks.NewChatServerRepositoryMock(mc)
				mock.ChatExistsMock.Expect(ctx, chatId).Return(true, nil)

				mock.DeleteBindUserFromChatMock.Expect(ctx, chatId).Return(query, nil)
				mock.CreateLogMock.Expect(ctx, log).Return(nil)

				mock.DeleteMessageFromChatMock.Expect(ctx, chatId).Return(query, nil)
				mock.CreateLogMock.Expect(ctx, log).Return(nil)

				mock.DeleteChatMock.Expect(ctx, chatId).Return(query, nil)
				mock.CreateLogMock.Expect(ctx, log).Return(nil)

				return mock
			},
		},
		{
			name: "error case ChatExistsMock",
			args: args{
				ctx:    ctx,
				chatId: chatId,
			},
			err: repoErr,
			chatServerRepositoryMock: func(mc *minimock.Controller) repository.ChatServerRepository {
				mock := repoMocks.NewChatServerRepositoryMock(mc)
				mock.ChatExistsMock.Expect(ctx, chatId).Return(true, repoErr)

				return mock
			},
		},
		{
			name: "error case not ChatExistsMock",
			args: args{
				ctx:    ctx,
				chatId: chatId,
			},
			err: errors.New("chat does not exist"),
			chatServerRepositoryMock: func(mc *minimock.Controller) repository.ChatServerRepository {
				mock := repoMocks.NewChatServerRepositoryMock(mc)
				mock.ChatExistsMock.Expect(ctx, chatId).Return(false, nil)

				return mock
			},
		},
		{
			name: "error case DeleteBindUserFromChatMock",
			args: args{
				ctx:    ctx,
				chatId: chatId,
			},
			err: repoErr,
			chatServerRepositoryMock: func(mc *minimock.Controller) repository.ChatServerRepository {
				mock := repoMocks.NewChatServerRepositoryMock(mc)
				mock.ChatExistsMock.Expect(ctx, chatId).Return(true, nil)

				mock.DeleteBindUserFromChatMock.Expect(ctx, chatId).Return(query, repoErr)
				return mock
			},
		},
		{
			name: "error case DeleteMessageFromChatMock",
			args: args{
				ctx:    ctx,
				chatId: chatId,
			},
			err: repoErr,
			chatServerRepositoryMock: func(mc *minimock.Controller) repository.ChatServerRepository {
				mock := repoMocks.NewChatServerRepositoryMock(mc)
				mock.ChatExistsMock.Expect(ctx, chatId).Return(true, nil)

				mock.DeleteBindUserFromChatMock.Expect(ctx, chatId).Return(query, nil)
				mock.CreateLogMock.Expect(ctx, log).Return(nil)

				mock.DeleteMessageFromChatMock.Expect(ctx, chatId).Return(query, nil)
				mock.CreateLogMock.Expect(ctx, log).Return(nil)

				mock.DeleteChatMock.Expect(ctx, chatId).Return(query, repoErr)
				return mock
			},
		},
		{
			name: "error case DeleteMessageFromChatMock",
			args: args{
				ctx:    ctx,
				chatId: chatId,
			},
			err: repoErr,
			chatServerRepositoryMock: func(mc *minimock.Controller) repository.ChatServerRepository {
				mock := repoMocks.NewChatServerRepositoryMock(mc)
				mock.ChatExistsMock.Expect(ctx, chatId).Return(true, nil)

				mock.DeleteBindUserFromChatMock.Expect(ctx, chatId).Return(query, nil)
				mock.CreateLogMock.Expect(ctx, log).Return(nil)

				mock.DeleteMessageFromChatMock.Expect(ctx, chatId).Return(query, repoErr)
				return mock
			},
		},

		{
			name: "create log error case after DeleteBindUserFromChatMock",
			args: args{
				ctx:    ctx,
				chatId: chatId,
			},
			err: repoErrLog,
			chatServerRepositoryMock: func(mc *minimock.Controller) repository.ChatServerRepository {
				mock := repoMocks.NewChatServerRepositoryMock(mc)
				mock.ChatExistsMock.Expect(ctx, chatId).Return(true, nil)

				mock.DeleteBindUserFromChatMock.Expect(ctx, chatId).Return(query, nil)
				mock.CreateLogMock.Expect(ctx, log).Return(repoErrLog)
				return mock
			},
		},
		//{
		//	name: "create log error case after DeleteChatMock",
		//	args: args{
		//		ctx:    ctx,
		//		chatId: chatId,
		//	},
		//	err: repoErrLog,
		//	chatServerRepositoryMock: func(mc *minimock.Controller) repository.ChatServerRepository {
		//		mock := repoMocks.NewChatServerRepositoryMock(mc)
		//		mock.ChatExistsMock.Expect(ctx, chatId).Return(true, nil)
		//
		//		mock.DeleteBindUserFromChatMock.Expect(ctx, chatId).Return(query, nil)
		//		mock.CreateLogMock.Expect(ctx, log).Return(nil)
		//
		//		mock.DeleteMessageFromChatMock.Expect(ctx, chatId).Return(query, nil)
		//		mock.CreateLogMock.Expect(ctx, log).Return(nil)
		//
		//		mock.DeleteChatMock.Expect(ctx, chatId).Return(query, nil)
		//		mock.CreateLogMock.Expect(ctx, log).Return(repoErrLog)
		//		return mock
		//	},
		//},
		//{
		//	name: "create log error case after DeleteMessageFromChatMock",
		//	args: args{
		//		ctx:    ctx,
		//		chatId: chatId,
		//	},
		//	err: repoErrLog,
		//	chatServerRepositoryMock: func(mc *minimock.Controller) repository.ChatServerRepository {
		//		mock := repoMocks.NewChatServerRepositoryMock(mc)
		//		mock.ChatExistsMock.Expect(ctx, chatId).Return(true, nil)
		//
		//		mock.DeleteBindUserFromChatMock.Expect(ctx, chatId).Return(query, nil)
		//		mock.CreateLogMock.Expect(ctx, log).Return(nil)
		//
		//		mock.DeleteMessageFromChatMock.Expect(ctx, chatId).Return(query, nil)
		//		mock.CreateLogMock.Expect(ctx, log).Return(repoErrLog)
		//		return mock
		//	},
		//},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(
			tt.name, func(t *testing.T) {
				t.Parallel()
				chatServerRepoMock := tt.chatServerRepositoryMock(mc)

				service := chatServerService.NewService(chatServerRepoMock, txManagerMock(), authServiceClientMock())

				err := service.DeleteChat(tt.args.ctx, tt.args.chatId)

				require.Equal(t, tt.err, err)
			},
		)
	}
}
