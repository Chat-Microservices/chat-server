package chatServer

import (
	"context"
	"github.com/semho/chat-microservices/chat-server/internal/client/db"
	"github.com/semho/chat-microservices/chat-server/internal/repository"
	"testing"
)

func Test_serv_CreateChat(t *testing.T) {
	type fields struct {
		chatServerRepository repository.ChatServerRepository
		txManager            db.TxManager
	}
	type args struct {
		ctx   context.Context
		users []*model.User
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				s := &serv{
					chatServerRepository: tt.fields.chatServerRepository,
					txManager:            tt.fields.txManager,
				}
				got, err := s.CreateChat(tt.args.ctx, tt.args.users)
				if (err != nil) != tt.wantErr {
					t.Errorf("CreateChat() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if got != tt.want {
					t.Errorf("CreateChat() got = %v, want %v", got, tt.want)
				}
			},
		)
	}
}
