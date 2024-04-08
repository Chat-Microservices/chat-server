package service

import (
	"context"
	"github.com/semho/chat-microservices/chat-server/internal/model"
)

type ChatServerService interface {
	CreateChat(ctx context.Context, users []*model.User) (int64, error)
	DeleteChat(ctx context.Context, chatId int64) error
	SendMessage(ctx context.Context, message *model.Message) error
	GetListLogs(ctx context.Context, pageNumber uint64, pageSize uint64) ([]*model.Log, error)
}
