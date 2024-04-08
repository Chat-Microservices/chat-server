package db

import (
	"context"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

// для работы с БД
type Client interface {
	DB() DB
	Close() error
}

// для работы с транзакциями
type Transactor interface {
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
}

// менеджер транзакций. выподняет обработчик указанный пользователем при транзакции
type TxManager interface {
	ReadCommitted(ctx context.Context, f Handler) error
}

// функция в транзакции
type Handler func(ctx context.Context) error

// обертка для запросов, хранит в себе название запроса и его SQL
// название используется для логирования
type Query struct {
	Name     string
	QueryRow string
}

// собирает NamedExecer и QueryExecer
type SQLExecer interface {
	NamedExecer
	QueryExecer
}

// для работы с именованными запросами с помощью тегов в структуре модели
type NamedExecer interface {
	ScanOneContext(ctx context.Context, dest any, q Query, args ...any) error
	ScanAllContext(ctx context.Context, dest any, q Query, args ...any) error
}

// для работы с обычными запросами
type QueryExecer interface {
	ExecContext(ctx context.Context, q Query, args ...any) (pgconn.CommandTag, error)
	QueryContext(ctx context.Context, q Query, args ...any) (pgx.Rows, error)
	QueryRowContext(ctx context.Context, q Query, args ...any) pgx.Row
}

// проверка соединения с БД
type Pinger interface {
	Ping(ctx context.Context) error
}

// работа с БД
type DB interface {
	SQLExecer
	Transactor
	Pinger
	Close()
}
