package converter

import (
	"github.com/semho/chat-microservices/chat-server/internal/model"
	desc "github.com/semho/chat-microservices/chat-server/pkg/chat-server_v1"
)

func ToMessageModelFromMessageApi(reqMessage *desc.SendMessageRequest) *model.Message {
	timestamp := reqMessage.Timestamp.AsTime()

	return &model.Message{
		UserName:  reqMessage.From,
		ChatID:    reqMessage.Id,
		Text:      reqMessage.Text,
		Timestamp: timestamp,
	}
}
