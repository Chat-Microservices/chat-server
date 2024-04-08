package modelRepo

import "time"

type Log struct {
	ID        int64     `db:"id"`
	Action    string    `db:"action"`
	EntityId  int64     `db:"entity_id"`
	Query     string    `db:"query"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
