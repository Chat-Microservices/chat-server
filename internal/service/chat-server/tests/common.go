package tests

import (
	"context"
	"github.com/semho/chat-microservices/chat-server/internal/client/db"
	"github.com/semho/chat-microservices/chat-server/internal/model"
)

type mockTxManager struct{}

func (m *mockTxManager) ReadCommitted(ctx context.Context, f db.Handler) error {
	return f(ctx)
}

type mockAuthServiceClient struct{}

func (m *mockAuthServiceClient) GetName(ctx context.Context, id int64) (*model.User, error) {
	return &model.User{Name: "MockName"}, nil
}
