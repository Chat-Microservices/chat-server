package model

import "time"

type Log struct {
	ID        int64
	Action    string
	EntityId  int64
	Query     string
	CreatedAt time.Time
	UpdatedAt time.Time
}
