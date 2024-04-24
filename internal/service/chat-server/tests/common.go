package tests

import (
	"context"
	"github.com/semho/chat-microservices/auth/internal/client/db"
)

type mockTxManager struct{}

func (m *mockTxManager) ReadCommitted(ctx context.Context, f db.Handler) error {
	return f(ctx)
}
