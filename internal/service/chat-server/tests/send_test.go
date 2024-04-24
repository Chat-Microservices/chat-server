package chatServer

import (
	"context"
	"github.com/semho/chat-microservices/chat-server/internal/client/db"
	"github.com/semho/chat-microservices/chat-server/internal/repository"
	"testing"
)

func Test_serv_SendMessage(t *testing.T) {
	type fields struct {
		chatServerRepository repository.ChatServerRepository
		txManager            db.TxManager
	}
	type args struct {
		ctx     context.Context
		message *model.Message
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
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
				if err := s.SendMessage(tt.args.ctx, tt.args.message); (err != nil) != tt.wantErr {
					t.Errorf("SendMessage() error = %v, wantErr %v", err, tt.wantErr)
				}
			},
		)
	}
}
