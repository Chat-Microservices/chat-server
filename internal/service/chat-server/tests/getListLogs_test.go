package chatServer

import (
	"context"
	"github.com/semho/chat-microservices/chat-server/internal/client/db"
	"github.com/semho/chat-microservices/chat-server/internal/repository"
	"reflect"
	"testing"
)

func Test_serv_GetListLogs(t *testing.T) {
	type fields struct {
		chatServerRepository repository.ChatServerRepository
		txManager            db.TxManager
	}
	type args struct {
		ctx        context.Context
		pageNumber uint64
		pageSize   uint64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*model.Log
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
				got, err := s.GetListLogs(tt.args.ctx, tt.args.pageNumber, tt.args.pageSize)
				if (err != nil) != tt.wantErr {
					t.Errorf("GetListLogs() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("GetListLogs() got = %v, want %v", got, tt.want)
				}
			},
		)
	}
}
