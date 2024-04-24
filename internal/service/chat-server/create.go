package chatServer

import (
	"context"
	"github.com/semho/chat-microservices/chat-server/internal/client/db"
	"github.com/semho/chat-microservices/chat-server/internal/converter"
	"github.com/semho/chat-microservices/chat-server/internal/model"
)

func (s *serv) CreateChat(ctx context.Context, users []*model.User) (int64, error) {

	var chatID int64
	var arrQuery []db.Query
	var q db.Query
	err := s.txManager.ReadCommitted(
		ctx, func(ctx context.Context) error {
			var errTx error
			var usersID []int64
			usersID, arrQuery, errTx = s.chatServerRepository.GetOfCreateUsers(ctx, users)
			if errTx != nil {
				return errTx
			}

			for index, id := range usersID {
				errTx = s.chatServerRepository.CreateLog(ctx, converter.ToChatServerLogFromQuery(arrQuery[index], id))
				if errTx != nil {
					return errTx
				}
			}

			chatID, q, errTx = s.chatServerRepository.CreateChat(ctx)
			if errTx != nil {
				return errTx
			}
			errTx = s.chatServerRepository.CreateLog(ctx, converter.ToChatServerLogFromQuery(q, chatID))
			if errTx != nil {
				return errTx
			}

			arrQuery, errTx = s.chatServerRepository.BindUserToChat(ctx, usersID, chatID)
			if errTx != nil {
				return errTx
			}
			for index, id := range usersID {
				errTx = s.chatServerRepository.CreateLog(ctx, converter.ToChatServerLogFromQuery(arrQuery[index], id))
				if errTx != nil {
					return errTx
				}
			}

			return nil
		},
	)

	if err != nil {
		return 0, err
	}

	return chatID, nil
}
