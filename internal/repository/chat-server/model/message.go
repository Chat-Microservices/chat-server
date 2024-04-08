package modelRepo

import "time"

type Message struct {
	ID        int64     `db:"id"`
	UserID    int64     `db:"user_id"`
	ChatID    int64     `db:"chat_id"`
	Text      string    `db:"text"`
	Timestamp time.Time `db:"timestamp"`
}
