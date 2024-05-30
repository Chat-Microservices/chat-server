package rpc

import (
	"context"
	"github.com/semho/chat-microservices/chat-server/internal/model"
)

type AuthServiceClient interface {
	GetName(ctx context.Context, id int64) (*model.User, error)
}
