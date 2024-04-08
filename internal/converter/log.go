package converter

import (
	"github.com/semho/chat-microservices/chat-server/internal/client/db"
	"github.com/semho/chat-microservices/chat-server/internal/model"
	desc "github.com/semho/chat-microservices/chat-server/pkg/chat-server_v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ToChatServerLogFromQuery(q db.Query, id int64) *model.Log {
	return &model.Log{
		Action:   q.Name,
		EntityId: id,
		Query:    q.QueryRow,
	}
}

func ToLogFromService(logs []*model.Log) *desc.LogsResponse {
	logsResponses := make([]*desc.Log, len(logs))
	for i, log := range logs {
		logsResponses[i] = &desc.Log{
			Id:        log.ID,
			Action:    log.Action,
			EntityId:  log.EntityId,
			Query:     log.Query,
			CreatedAt: timestamppb.New(log.CreatedAt),
			UpdatedAt: timestamppb.New(log.UpdatedAt),
		}

	}
	return &desc.LogsResponse{
		Logs: logsResponses,
	}
}
