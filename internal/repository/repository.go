package repository

import (
	"context"
	"github.com/semho/chat-microservices/chat-server/internal/client/db"
	"github.com/semho/chat-microservices/chat-server/internal/model"
)

type ChatServerRepository interface {
	GetUserIDByName(ctx context.Context, name string) (int64, db.Query, error)
	GetOfCreateUsers(ctx context.Context, users []*model.User) ([]int64, []db.Query, error)
	CreateChat(ctx context.Context) (int64, db.Query, error)
	BindUserToChat(ctx context.Context, usersID []int64, chatID int64) ([]db.Query, error)

	DeleteBindUserFromChat(ctx context.Context, chatID int64) (db.Query, error)
	DeleteMessageFromChat(ctx context.Context, chatID int64) (db.Query, error)
	DeleteChat(ctx context.Context, chatID int64) (db.Query, error)

	ChatExists(ctx context.Context, chatID int64) (bool, error)

	SendMessage(ctx context.Context, message *model.Message, userID int64) (db.Query, error)

	CreateLog(ctx context.Context, logger *model.Log) error
	GetListLog(ctx context.Context, pageNumber uint64, pageSize uint64) ([]*model.Log, error)
}
