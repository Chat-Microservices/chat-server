package chatServerAPI

import (
	"context"
	"github.com/semho/chat-microservices/chat-server/internal/service"
	chat_server_v1 "github.com/semho/chat-microservices/chat-server/pkg/chat-server_v1"
	"reflect"
	"testing"
)

func TestImplementation_GetListLogs(t *testing.T) {
	type fields struct {
		UnimplementedChatServerV1Server chat_server_v1.UnimplementedChatServerV1Server
		chatServerService               service.ChatServerService
	}
	type args struct {
		ctx context.Context
		req *desc.GetListLogsRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *desc.LogsResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				i := &Implementation{
					UnimplementedChatServerV1Server: tt.fields.UnimplementedChatServerV1Server,
					chatServerService:               tt.fields.chatServerService,
				}
				got, err := i.GetListLogs(tt.args.ctx, tt.args.req)
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
