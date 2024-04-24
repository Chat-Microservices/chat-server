package chatServer

import (
	"context"
	"errors"
	"github.com/semho/chat-microservices/chat-server/internal/client/db"
	"github.com/semho/chat-microservices/chat-server/internal/converter"
)

func (s *serv) DeleteChat(ctx context.Context, chatId int64) error {
	var ok bool
	var err error

	ok, err = s.chatServerRepository.ChatExists(ctx, chatId)
	if err != nil {
		return err
	}

	if !ok {
		return errors.New("chat does not exist")
	}

	err = s.txManager.ReadCommitted(
		ctx, func(ctx context.Context) error {
			var errTx error
			var q db.Query
			q, errTx = s.chatServerRepository.DeleteBindUserFromChat(ctx, chatId)
			if errTx != nil {
				return errTx
			}
			errTx = s.chatServerRepository.CreateLog(ctx, converter.ToChatServerLogFromQuery(q, chatId))
			if errTx != nil {
				return errTx
			}

			q, errTx = s.chatServerRepository.DeleteMessageFromChat(ctx, chatId)
			if errTx != nil {
				return errTx
			}
			errTx = s.chatServerRepository.CreateLog(ctx, converter.ToChatServerLogFromQuery(q, chatId))
			if errTx != nil {
				return errTx
			}

			q, errTx = s.chatServerRepository.DeleteChat(ctx, chatId)
			if errTx != nil {
				return errTx
			}
			errTx = s.chatServerRepository.CreateLog(ctx, converter.ToChatServerLogFromQuery(q, chatId))
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
