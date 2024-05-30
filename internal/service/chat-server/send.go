package chatServer

import (
	"context"
	"errors"
	"github.com/fatih/color"
	"github.com/opentracing/opentracing-go"
	"github.com/semho/chat-microservices/chat-server/internal/client/db"
	"github.com/semho/chat-microservices/chat-server/internal/converter"
	"github.com/semho/chat-microservices/chat-server/internal/model"
	"log"
)

const (
	userID = 1
)

func (s *serv) SendMessage(ctx context.Context, message *model.Message) error {
	var ok bool
	var err error

	ok, err = s.chatServerRepository.ChatExists(ctx, message.ChatID)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("chat does not exist")
	}

	//TODO: как пример отправки запроса клиента с трейсингом
	span, ctx := opentracing.StartSpanFromContext(ctx, "get user")
	defer span.Finish()

	span.SetTag("id", userID)
	user, err := s.authServiceClient.GetName(ctx, userID)
	if err != nil {
		return err
	}

	log.Printf(color.RedString("Answer: \n"), color.GreenString("%+v", user.Name))

	var q db.Query
	var userID int64
	err = s.txManager.ReadCommitted(
		ctx, func(ctx context.Context) error {
			var errTx error

			userID, q, errTx = s.chatServerRepository.GetUserIDByName(ctx, message.UserName)
			if errTx != nil {
				return errTx
			}
			errTx = s.chatServerRepository.CreateLog(ctx, converter.ToChatServerLogFromQuery(q, userID))
			if errTx != nil {
				return errTx
			}

			q, errTx = s.chatServerRepository.SendMessage(ctx, message, userID)
			if errTx != nil {
				return errTx
			}
			errTx = s.chatServerRepository.CreateLog(ctx, converter.ToChatServerLogFromQuery(q, userID))
			if errTx != nil {
				return errTx
			}

			return nil
		},
	)

	if err != nil {
		return err
	}

	return nil
}
