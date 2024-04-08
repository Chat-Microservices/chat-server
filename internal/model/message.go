package model

import "time"

type Message struct {
	UserName  string
	ChatID    int64
	Text      string
	Timestamp time.Time
}
