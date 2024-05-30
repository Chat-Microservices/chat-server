package auth_client

import (
	"context"
	"github.com/semho/chat-microservices/chat-server/internal/model"
	desc "github.com/semho/chat-microservices/chat-server/pkg/auth_v1"
)

type client struct {
	authClient desc.AuthV1Client
}

func New(authClient desc.AuthV1Client) *client {
	return &client{authClient: authClient}
}

func (c *client) GetName(ctx context.Context, id int64) (*model.User, error) {
	res, err := c.authClient.Get(ctx, &desc.GetRequest{Id: id})
	if err != nil {
		return nil, err
	}

	return &model.User{Name: res.GetDetail().GetName()}, nil
}
