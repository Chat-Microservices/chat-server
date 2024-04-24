package tests

import (
	"context"
	"fmt"
	"github.com/brianvoe/gofakeit"
	"github.com/gojuno/minimock/v3"
	"github.com/semho/chat-microservices/chat-server/internal/client/db"
	"github.com/semho/chat-microservices/chat-server/internal/model"
	"github.com/semho/chat-microservices/chat-server/internal/repository"
	repoMocks "github.com/semho/chat-microservices/chat-server/internal/repository/mocks"
	chatServerService "github.com/semho/chat-microservices/chat-server/internal/service/chat-server"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_serv_CreateChat(t *testing.T) {
	t.Parallel()
	type chatServerRepositoryMockFunc func(mc *minimock.Controller) repository.ChatServerRepository

	type args struct {
		ctx   context.Context
		users []*model.User
	}

	txManagerMock := func() db.TxManager {
		return &mockTxManager{}
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		name1 = gofakeit.Name()
		name2 = gofakeit.Name()
		users = []*model.User{
			{
				Name: name1,
			},
			{
				Name: name2,
			},
		}
		id = gofakeit.Int64()

		qName = "CreateLog"
		qRow  = "INSERT INTO..."
		query = db.Query{
			Name:     qName,
			QueryRow: qRow,
		}

		queryArr = []db.Query{
			query, query,
		}

		usersId    = []int64{id, id}
		repoErr    = fmt.Errorf("repo error")
		repoErrLog = fmt.Errorf("repo error log")
		log        = &model.Log{
			Action:   qName,
			EntityId: id,
			Query:    qRow,
		}
	)

	tests := []struct {
		name                     string
		args                     args
		want                     int64
		err                      error
		chatServerRepositoryMock chatServerRepositoryMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx:   ctx,
				users: users,
			},
			want: id,
			err:  nil,
			chatServerRepositoryMock: func(mc *minimock.Controller) repository.ChatServerRepository {
				mock := repoMocks.NewChatServerRepositoryMock(mc)
				mock.GetOfCreateUsersMock.Expect(ctx, users).Return(usersId, queryArr, nil)
				for _, userId := range usersId {
					mock.CreateLogMock.Expect(
						ctx, &model.Log{
							Action:   qName,
							EntityId: userId,
							Query:    qRow,
						},
					).Return(nil)
				}

				mock.CreateChatMock.Expect(ctx).Return(id, query, nil)
				mock.CreateLogMock.Expect(ctx, log).Return(nil)

				mock.BindUserToChatMock.Expect(ctx, usersId, id).Return(queryArr, nil)
				for _, userId := range usersId {
					mock.CreateLogMock.Expect(
						ctx, &model.Log{
							Action:   qName,
							EntityId: userId,
							Query:    qRow,
						},
					).Return(nil)
				}
				return mock
			},
		},
		{
			name: "error GetOfCreateUsersMock",
			args: args{
				ctx:   ctx,
				users: users,
			},
			want: 0,
			err:  repoErr,
			chatServerRepositoryMock: func(mc *minimock.Controller) repository.ChatServerRepository {
				mock := repoMocks.NewChatServerRepositoryMock(mc)
				mock.GetOfCreateUsersMock.Expect(ctx, users).Return(nil, queryArr, repoErr)
				return mock
			},
		},
		{
			name: "error CreateChatMock",
			args: args{
				ctx:   ctx,
				users: users,
			},
			want: 0,
			err:  repoErr,
			chatServerRepositoryMock: func(mc *minimock.Controller) repository.ChatServerRepository {
				mock := repoMocks.NewChatServerRepositoryMock(mc)
				mock.GetOfCreateUsersMock.Expect(ctx, users).Return(usersId, queryArr, nil)
				for _, userId := range usersId {
					mock.CreateLogMock.Expect(
						ctx, &model.Log{
							Action:   qName,
							EntityId: userId,
							Query:    qRow,
						},
					).Return(nil)
				}
				mock.CreateChatMock.Expect(ctx).Return(0, query, repoErr)
				return mock
			},
		},
		{
			name: "error BindUserToChatMock",
			args: args{
				ctx:   ctx,
				users: users,
			},
			want: 0,
			err:  repoErr,
			chatServerRepositoryMock: func(mc *minimock.Controller) repository.ChatServerRepository {
				mock := repoMocks.NewChatServerRepositoryMock(mc)
				mock.GetOfCreateUsersMock.Expect(ctx, users).Return(usersId, queryArr, nil)
				for _, userId := range usersId {
					mock.CreateLogMock.Expect(
						ctx, &model.Log{
							Action:   qName,
							EntityId: userId,
							Query:    qRow,
						},
					).Return(nil)
				}
				mock.CreateChatMock.Expect(ctx).Return(id, query, nil)
				mock.CreateLogMock.Expect(ctx, log).Return(nil)

				mock.BindUserToChatMock.Expect(ctx, usersId, id).Return(nil, repoErr)

				return mock
			},
		},
		{
			name: "create log error case GetOfCreateUsersMock",
			args: args{
				ctx:   ctx,
				users: users,
			},
			want: 0,
			err:  repoErrLog,
			chatServerRepositoryMock: func(mc *minimock.Controller) repository.ChatServerRepository {
				mock := repoMocks.NewChatServerRepositoryMock(mc)
				mock.GetOfCreateUsersMock.Expect(ctx, users).Return(usersId, queryArr, nil)
				mock.CreateLogMock.Expect(ctx, log).Return(repoErrLog)
				return mock
			},
		},
		//{
		//	name: "create log error case CreateChatMock",
		//	args: args{
		//		ctx:   ctx,
		//		users: users,
		//	},
		//	want: 0,
		//	err:  repoErrLog,
		//	chatServerRepositoryMock: func(mc *minimock.Controller) repository.ChatServerRepository {
		//		mock := repoMocks.NewChatServerRepositoryMock(mc)
		//		mock.GetOfCreateUsersMock.Expect(ctx, users).Return(usersId, queryArr, nil)
		//
		//		for _, userId := range usersId {
		//			mock.CreateLogMock.Expect(
		//				ctx, &model.Log{
		//					Action:   qName,
		//					EntityId: userId,
		//					Query:    qRow,
		//				},
		//			).Return(nil)
		//		}
		//		mock.CreateChatMock.Expect(ctx).Return(id, query, nil)
		//		mock.CreateLogMock.Expect(ctx, log).Return(repoErrLog)
		//
		//		return mock
		//	},
		//},
		//{
		//	name: "create log error case BindUserToChatMock",
		//	args: args{
		//		ctx:   ctx,
		//		users: users,
		//	},
		//	want: 0,
		//	err:  repoErrLog,
		//	chatServerRepositoryMock: func(mc *minimock.Controller) repository.ChatServerRepository {
		//		mock := repoMocks.NewChatServerRepositoryMock(mc)
		//		mock.GetOfCreateUsersMock.Expect(ctx, users).Return(usersId, queryArr, nil)
		//		for _, userId := range usersId {
		//			mock.CreateLogMock.Expect(
		//				ctx, &model.Log{
		//					Action:   qName,
		//					EntityId: userId,
		//					Query:    qRow,
		//				},
		//			).Return(nil)
		//		}
		//
		//		mock.CreateChatMock.Expect(ctx).Return(id, query, nil)
		//		mock.CreateLogMock.Expect(
		//			ctx, &model.Log{
		//				Action:   qName,
		//				EntityId: id,
		//				Query:    qRow,
		//			},
		//		).Return(nil)
		//
		//		mock.BindUserToChatMock.Expect(ctx, usersId, id).Return(queryArr, nil)
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

				service := chatServerService.NewService(chatServerRepoMock, txManagerMock())

				resHandler, err := service.CreateChat(tt.args.ctx, tt.args.users)

				require.Equal(t, tt.err, err)
				require.Equal(t, tt.want, resHandler)
			},
		)
	}
}
