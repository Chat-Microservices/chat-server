package main

import (
	"context"
	"database/sql"
	"flag"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/semho/chat-microservices/chat-server/internal/config"
	"github.com/semho/chat-microservices/chat-server/internal/config/env"
	desc "github.com/semho/chat-microservices/chat-server/pkg/chat-server_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"net"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", ".env", "path to config file")
}

type server struct {
	desc.UnimplementedChatServerV1Server
	pool *pgxpool.Pool
}

var internalServerError = status.Error(codes.Internal, "Internal server error")

func checkError(msg string, err error) error {
	log.Printf("%s: %v", msg, err)
	return internalServerError
}

func (s *server) userExists(ctx context.Context, userName string) (bool, error) {
	var exists bool
	err := s.pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM "users"  WHERE name = $1)`,
		userName).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (s *server) chatExists(ctx context.Context, chatID int64) (bool, error) {
	var exists bool
	err := s.pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM "chats"  WHERE id = $1)`,
		chatID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (s *server) CreateChat(ctx context.Context, req *desc.CreateChatRequest) (*desc.CreateChatResponse, error) {
	if req.GetUsernames() == nil {
		return nil, status.Error(codes.InvalidArgument, "Invalid request: UserNames must be provided")
	}
	log.Printf("User names: %+v", req.GetUsernames())

	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel:       pgx.ReadCommitted,
		AccessMode:     pgx.ReadWrite,
		DeferrableMode: pgx.NotDeferrable,
	})
	if err != nil {
		return nil, checkError("Error starting transaction", err)
	}
	defer tx.Rollback(ctx)

	selectUserBuilder := sq.Select("id").From("users").PlaceholderFormat(sq.Dollar)
	insertUserBuilder := sq.Insert("users").PlaceholderFormat(sq.Dollar).Columns("name").Suffix("RETURNING id")
	insertUserChatBuilder := sq.Insert("user_chat").PlaceholderFormat(sq.Dollar).Columns("user_id", "chat_id")
	usersID := make([]int64, len(req.GetUsernames()))
	for i, user := range req.GetUsernames() {
		exists, err := s.userExists(ctx, user.GetName())
		if err != nil {
			return nil, checkError("Error checking user existence", err)
		}
		if exists {
			query, args, err := selectUserBuilder.Where(sq.Eq{"name": user.GetName()}).ToSql()
			if err != nil {
				return nil, checkError("Failed to build query", err)
			}
			err = tx.QueryRow(ctx, query, args...).Scan(&usersID[i])
			if err != nil {
				return nil, checkError("Failed to select user from the database", err)
			}
		} else {
			query, args, err := insertUserBuilder.Values(user.GetName()).ToSql()
			if err != nil {
				return nil, checkError("Failed to build query", err)
			}

			err = tx.QueryRow(ctx, query, args...).Scan(&usersID[i])
			if err != nil {
				return nil, checkError("Failed to insert user into the database", err)
			}
		}
	}

	var chatID int64
	err = tx.QueryRow(ctx, `INSERT INTO chats DEFAULT VALUES RETURNING id`).Scan(&chatID)
	if err != nil {
		return nil, checkError("failed to create chat the database", err)
	}

	for _, userID := range usersID {
		query, args, err := insertUserChatBuilder.Values(userID, chatID).ToSql()
		if err != nil {
			return nil, checkError("Failed to build query", err)
		}
		_, err = tx.Exec(ctx, query, args...)
		if err != nil {
			return nil, checkError("failed to insert userID and chatID into the database", err)
		}
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, checkError("Error committing transaction", err)
	}

	newChat := &desc.Chat{
		Id:        chatID,
		Usernames: req.GetUsernames(),
	}
	log.Printf("New chat: %+v", newChat)

	response := &desc.CreateChatResponse{
		Id: newChat.Id,
	}

	return response, nil
}

func (s *server) DeleteChat(ctx context.Context, req *desc.DeleteRequest) (*emptypb.Empty, error) {
	if req.GetId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "Invalid request: Id must be provided")
	}

	exists, err := s.chatExists(ctx, req.GetId())
	if err != nil {
		return nil, checkError("Error checking chat existence", err)
	}
	if !exists {
		return nil, status.Error(codes.NotFound, "Chat not found")
	}

	log.Printf("delete chat by id: %d", req.GetId())

	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel:       pgx.ReadCommitted,
		AccessMode:     pgx.ReadWrite,
		DeferrableMode: pgx.NotDeferrable,
	})
	if err != nil {
		return nil, checkError("Error starting transaction", err)
	}
	defer tx.Rollback(ctx)

	deleteUserChatBuilder := sq.Delete("user_chat").PlaceholderFormat(sq.Dollar).Where(sq.Eq{"chat_id": req.GetId()})
	deleteMessagesBuilder := sq.Delete("messages").PlaceholderFormat(sq.Dollar).Where(sq.Eq{"chat_id": req.GetId()})
	deleteChatBuilder := sq.Delete("chats").PlaceholderFormat(sq.Dollar).Where(sq.Eq{"id": req.GetId()})
	sqlUserChat, argsUserChat, err := deleteUserChatBuilder.ToSql()
	if err != nil {
		return nil, checkError("Failed to generate SQL for deleting user_chat records", err)
	}
	sqlMessages, argsMessages, err := deleteMessagesBuilder.ToSql()
	if err != nil {
		return nil, checkError("Failed to generate SQL for deleting messages records", err)
	}
	sqlChat, argsChat, err := deleteChatBuilder.ToSql()
	if err != nil {
		return nil, checkError("Failed to generate SQL for deleting chat", err)
	}
	_, err = tx.Exec(ctx, sqlUserChat, argsUserChat...)
	if err != nil {
		return nil, checkError("failed to delete user_chat records from the database", err)
	}

	_, err = tx.Exec(ctx, sqlMessages, argsMessages...)
	if err != nil {
		return nil, checkError("failed to delete user_chat records from the database", err)
	}

	_, err = tx.Exec(ctx, sqlChat, argsChat...)
	if err != nil {
		return nil, checkError("failed to delete chat the database", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, checkError("Error committing transaction", err)
	}

	return &emptypb.Empty{}, nil
}

func (s *server) SendMessage(ctx context.Context, req *desc.SendMessageRequest) (*emptypb.Empty, error) {
	if req.GetFrom() == "" || req.GetText() == "" {
		return nil, status.Error(codes.InvalidArgument, "Invalid request: From or Text must be provided")
	}
	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel:       pgx.ReadCommitted,
		AccessMode:     pgx.ReadWrite,
		DeferrableMode: pgx.NotDeferrable,
	})
	if err != nil {
		return nil, checkError("Error starting transaction", err)
	}
	defer tx.Rollback(ctx)

	selectUserIDQuery, selectUserIDArgs, err := sq.
		Select("id").
		From("users").
		Where(sq.Eq{"name": req.GetFrom()}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, checkError("Failed to build select user query", err)
	}

	var userID int64
	err = tx.QueryRow(ctx, selectUserIDQuery, selectUserIDArgs...).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Error(codes.NotFound, "User not found")
		}
		return nil, checkError("Failed to select user from the database", err)
	}

	selectChatIDQuery, selectChatIDArgs, err := sq.
		Select("DISTINCT chat_id").
		From("user_chat").
		Where(sq.Eq{"user_id": userID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, checkError("Failed to build select chatID query", err)
	}

	var chatID int64
	err = tx.QueryRow(ctx, selectChatIDQuery, selectChatIDArgs...).Scan(&chatID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Error(codes.NotFound, "Chat not found")
		}
		return nil, checkError("Failed to select chatID from the database", err)
	}

	log.Printf("UserID: %d, ChatID: %d", userID, chatID)

	insertMessageQuery, insertMessageArgs, err := sq.Insert("messages").
		PlaceholderFormat(sq.Dollar).
		Columns("user_id", "chat_id", "text").
		Values(userID, chatID, req.GetText()).
		ToSql()
	if err != nil {
		return nil, checkError("Failed to build insert message query", err)
	}

	_, err = tx.Exec(ctx, insertMessageQuery, insertMessageArgs...)
	if err != nil {
		return nil, checkError("failed to create message in the database", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, checkError("Error committing transaction", err)
	}

	newMessage := &desc.SendMessageRequest{
		From: req.GetFrom(),
		Text: req.GetText(),
	}

	log.Printf("New message: %+v", newMessage)

	return &emptypb.Empty{}, nil
}

func main() {
	flag.Parse()
	ctx := context.Background()

	err := config.Load(configPath)
	if err != nil {
		log.Fatalf("failed to load config: %s", err)
	}

	grpcConfig, err := env.NewGRPCConfig()
	if err != nil {
		log.Fatalf("failed to get grpc config: %v", err)
	}

	pgConfig, err := env.NewPGConfig()
	if err != nil {
		log.Fatalf("failed to get pg config: %v", err)
	}

	lis, err := net.Listen("tcp", grpcConfig.Address())
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	pool, err := pgxpool.Connect(ctx, pgConfig.DSN())
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	defer pool.Close()

	s := grpc.NewServer()
	reflection.Register(s)
	desc.RegisterChatServerV1Server(s, &server{pool: pool})

	log.Printf("server listening at: %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
