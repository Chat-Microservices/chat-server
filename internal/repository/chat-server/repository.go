package chat_server

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/semho/chat-microservices/chat-server/internal/client/db"
	"github.com/semho/chat-microservices/chat-server/internal/model"
	"github.com/semho/chat-microservices/chat-server/internal/repository"
	"github.com/semho/chat-microservices/chat-server/internal/repository/chat-server/converter"
	modelRepo "github.com/semho/chat-microservices/chat-server/internal/repository/chat-server/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
)

const (
	tableUsers = "users"

	idColumn   = "id"
	nameColumn = "name"

	tableChats = "chats"

	tableMessages = "messages"

	userIDColumn    = "user_id"
	chatIDCColumn   = "chat_id"
	textColumn      = "text"
	timestampColumn = "timestamp"

	tableUserChat = "user_chat"

	tableLogs       = "logs"
	actionColumn    = "action"
	entityColumnID  = "entity_id"
	queryColumn     = "query"
	createdAtColumn = "created_at"
	updatedAtColumn = "updated_at"
)

type repo struct {
	db db.Client
}

func NewRepository(db db.Client) repository.ChatServerRepository {
	return &repo{db: db}
}

func (r *repo) GetUserIDByName(ctx context.Context, name string) (int64, db.Query, error) {
	_, err := r.userExists(ctx, name)
	if err != nil {
		return 0, db.Query{}, checkError("Error checking user existence", err)
	}

	selectUserBuilder := sq.Select(idColumn).From(tableUsers).PlaceholderFormat(sq.Dollar)
	query, args, err := selectUserBuilder.Where(sq.Eq{nameColumn: name}).ToSql()
	if err != nil {
		return 0, db.Query{}, checkError("Failed to build query", err)
	}
	q := db.Query{
		Name:     "user_repository.Get",
		QueryRow: query,
	}
	var usersID int64
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&usersID)
	if err != nil {
		return 0, db.Query{}, checkError("Failed to select user from the database", err)
	}

	return usersID, q, nil
}

func (r *repo) GetOfCreateUsers(ctx context.Context, users []*model.User) ([]int64, []db.Query, error) {
	insertUserBuilder := sq.Insert(tableUsers).PlaceholderFormat(sq.Dollar).Columns(nameColumn).Suffix("RETURNING id")

	arrQuery := make([]db.Query, len(users))
	var q db.Query
	usersID := make([]int64, len(users))
	for i, user := range users {
		exists, err := r.userExists(ctx, user.Name)
		if err != nil {
			return nil, []db.Query{}, checkError("Error checking user existence", err)
		}
		if exists {
			var userID int64
			userID, q, err = r.GetUserIDByName(ctx, user.Name)
			if err != nil {
				return nil, []db.Query{}, checkError("Failed to get user id", err)
			}
			usersID[i] = userID
		} else {
			query, args, err := insertUserBuilder.Values(user.Name).ToSql()
			if err != nil {
				return nil, []db.Query{}, checkError("Failed to build query", err)
			}

			q = db.Query{
				Name:     "user_repository.Create",
				QueryRow: query,
			}

			err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&usersID[i])
			if err != nil {
				return nil, []db.Query{}, checkError("Failed to insert user into the database", err)
			}
		}
		arrQuery[i] = q
	}

	return usersID, arrQuery, nil
}

func (r repo) CreateChat(ctx context.Context) (int64, db.Query, error) {
	var chatID int64
	q := db.Query{
		Name:     "chat_repository.Create",
		QueryRow: "INSERT INTO chats DEFAULT VALUES RETURNING id",
	}
	err := r.db.DB().QueryRowContext(ctx, q).Scan(&chatID)
	if err != nil {
		return 0, db.Query{}, checkError("failed to create chat the database", err)
	}

	return chatID, q, nil
}

func (r repo) BindUserToChat(ctx context.Context, usersID []int64, chatID int64) ([]db.Query, error) {
	insertUserChatBuilder := sq.Insert(tableUserChat).PlaceholderFormat(sq.Dollar).Columns(userIDColumn, chatIDCColumn)
	var q db.Query
	arrQuery := make([]db.Query, len(usersID))
	for i, userID := range usersID {
		query, args, err := insertUserChatBuilder.Values(userID, chatID).ToSql()
		if err != nil {
			return []db.Query{}, checkError("Failed to build query", err)
		}

		q = db.Query{
			Name:     "user_chat_repository.Create",
			QueryRow: query,
		}
		_, err = r.db.DB().ExecContext(ctx, q, args...)
		if err != nil {
			return []db.Query{}, checkError("failed to insert userID and chatID into the database", err)
		}

		arrQuery[i] = q
	}

	return arrQuery, nil
}

func (r *repo) DeleteBindUserFromChat(ctx context.Context, chatID int64) (db.Query, error) {
	deleteUserChatBuilder := sq.Delete(tableUserChat).PlaceholderFormat(sq.Dollar).Where(sq.Eq{chatIDCColumn: chatID})
	query, args, err := deleteUserChatBuilder.ToSql()
	if err != nil {
		return db.Query{}, checkError("Failed to generate SQL for deleting user_chat records", err)
	}

	q := db.Query{
		Name:     "user_chat_repository.Delete",
		QueryRow: query,
	}

	_, err = r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return db.Query{}, checkError("failed to delete user_chat records from the database", err)
	}

	return q, nil
}

func (r *repo) DeleteMessageFromChat(ctx context.Context, chatID int64) (db.Query, error) {
	deleteMessagesBuilder := sq.Delete(tableMessages).PlaceholderFormat(sq.Dollar).Where(sq.Eq{chatIDCColumn: chatID})
	query, args, err := deleteMessagesBuilder.ToSql()
	if err != nil {
		return db.Query{}, checkError("Failed to generate SQL for deleting messages records", err)
	}

	q := db.Query{
		Name:     "message_repository.Delete",
		QueryRow: query,
	}

	_, err = r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return db.Query{}, checkError("failed to delete user_chat records from the database", err)
	}

	return q, nil
}

func (r repo) DeleteChat(ctx context.Context, chatID int64) (db.Query, error) {
	deleteChatBuilder := sq.Delete(tableChats).PlaceholderFormat(sq.Dollar).Where(sq.Eq{idColumn: chatID})
	query, args, err := deleteChatBuilder.ToSql()
	if err != nil {
		return db.Query{}, checkError("Failed to generate SQL for deleting chat", err)
	}

	q := db.Query{
		Name:     "chat_repository.Delete",
		QueryRow: query,
	}

	_, err = r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return db.Query{}, checkError("failed to delete chat the database", err)
	}

	return q, nil
}

func (r repo) SendMessage(ctx context.Context, message *model.Message, userID int64) (db.Query, error) {
	query, arg, err := sq.Insert(tableMessages).
		PlaceholderFormat(sq.Dollar).
		Columns(userIDColumn, chatIDCColumn, textColumn).
		Values(userID, message.ChatID, message.Text).
		ToSql()
	if err != nil {
		return db.Query{}, checkError("Failed to build insert message query", err)
	}

	q := db.Query{
		Name:     "message_repository.Create",
		QueryRow: query,
	}

	_, err = r.db.DB().ExecContext(ctx, q, arg...)
	if err != nil {
		return db.Query{}, checkError("failed to create message in the database", err)
	}

	return q, nil
}

var internalServerError = status.Error(codes.Internal, "Internal server error")

func checkError(msg string, err error) error {
	log.Printf("%s: %v", msg, err)
	return internalServerError
}

func (r *repo) userExists(ctx context.Context, userName string) (bool, error) {
	var exists bool
	q := db.Query{
		Name:     "user_repository.Exist",
		QueryRow: "SELECT EXISTS(SELECT 1 FROM users WHERE name = $1)",
	}
	err := r.db.DB().QueryRowContext(ctx, q, userName).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (r *repo) ChatExists(ctx context.Context, chatID int64) (bool, error) {
	var exists bool
	q := db.Query{
		Name:     "chat_repository.Exist",
		QueryRow: "SELECT EXISTS(SELECT 1 FROM chats WHERE id = $1)",
	}
	err := r.db.DB().QueryRowContext(ctx, q, chatID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (r repo) CreateLog(ctx context.Context, logger *model.Log) error {
	query, args, err := sq.Insert(tableLogs).
		PlaceholderFormat(sq.Dollar).
		Columns(actionColumn, entityColumnID, queryColumn).
		Values(logger.Action, logger.EntityId, logger.Query).
		ToSql()

	if err != nil {
		log.Printf("failed to build query: %v", err)
		return status.Error(codes.Internal, "Internal server error")
	}

	q := db.Query{
		Name:     "log_repository.Create",
		QueryRow: query,
	}

	if _, err = r.db.DB().ExecContext(ctx, q, args...); err != nil {
		log.Printf("failed to create log: %v", err)
		return status.Error(codes.Internal, "Internal server error")
	}

	return nil
}

func (r repo) GetListLog(ctx context.Context, pageNumber uint64, pageSize uint64) ([]*model.Log, error) {
	offset := (pageNumber - 1) * pageSize

	query, args, err := sq.Select(idColumn, actionColumn, entityColumnID, queryColumn, createdAtColumn, updatedAtColumn).
		From(tableLogs).
		PlaceholderFormat(sq.Dollar).
		Limit(pageSize).
		Offset(offset).
		ToSql()
	if err != nil {
		log.Println(err)
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	q := db.Query{
		Name:     "log_repository.Get",
		QueryRow: query,
	}

	var logs []modelRepo.Log
	err = r.db.DB().ScanAllContext(ctx, &logs, q, args...)
	if err != nil {
		log.Println(err)
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	return converter.ToLogFromRepo(logs), nil
}
