package converter

import (
	"github.com/semho/chat-microservices/chat-server/internal/model"
	desc "github.com/semho/chat-microservices/chat-server/pkg/chat-server_v1"
)

func ToUserModelFromUserApi(api []*desc.User) []*model.User {
	users := make([]*model.User, len(api))

	for i, u := range api {
		users[i] = &model.User{
			Name: u.Name,
		}
	}

	return users
}
