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
	"time"
)

func Test_serv_GetListLogs(t *testing.T) {
	t.Parallel()
	type chatServerRepositoryMockFunc func(mc *minimock.Controller) repository.ChatServerRepository

	txManagerMock := func() db.TxManager {
		return &mockTxManager{}
	}
	type args struct {
		ctx        context.Context
		pageNumber uint64
		pageSize   uint64
	}

	var (
		ctx        = context.Background()
		mc         = minimock.NewController(t)
		pageNumber = gofakeit.Uint64()
		pageSize   = gofakeit.Uint64()
		created    = time.Now()
		res        = []*model.Log{
			{
				ID:        gofakeit.Int64(),
				Action:    gofakeit.BeerName(),
				EntityId:  gofakeit.Int64(),
				Query:     gofakeit.City(),
				CreatedAt: created,
				UpdatedAt: created,
			},
		}
		repoErr = fmt.Errorf("repo error")
	)

	tests := []struct {
		name                     string
		args                     args
		want                     []*model.Log
		err                      error
		chatServerRepositoryMock chatServerRepositoryMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx:        ctx,
				pageNumber: pageNumber,
				pageSize:   pageSize,
			},
			want: res,
			err:  nil,
			chatServerRepositoryMock: func(mc *minimock.Controller) repository.ChatServerRepository {
				mock := repoMocks.NewChatServerRepositoryMock(mc)
				mock.GetListLogMock.Expect(ctx, pageNumber, pageSize).Return(res, nil)
				return mock
			},
		},
		{
			name: "success case",
			args: args{
				ctx:        ctx,
				pageNumber: pageNumber,
				pageSize:   pageSize,
			},
			want: nil,
			err:  repoErr,
			chatServerRepositoryMock: func(mc *minimock.Controller) repository.ChatServerRepository {
				mock := repoMocks.NewChatServerRepositoryMock(mc)
				mock.GetListLogMock.Expect(ctx, pageNumber, pageSize).Return(nil, repoErr)
				return mock
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(
			tt.name, func(t *testing.T) {
				t.Parallel()
				chatServerRepoMock := tt.chatServerRepositoryMock(mc)

				service := chatServerService.NewService(chatServerRepoMock, txManagerMock())

				resHandler, err := service.GetListLogs(tt.args.ctx, tt.args.pageNumber, tt.args.pageSize)

				require.Equal(t, tt.err, err)
				require.Equal(t, tt.want, resHandler)
			},
		)
	}
}
