package models

import "time"

type Event struct {
	ID        int       `json:"id" bun:"id,pk,autoincrement"`
	TaskID    int       `json:"task_id" bun:"task_id,notnull"`
	Name      string    `json:"name" bun:"name,notnull"`
	CreatedAt time.Time `json:"created_at" bun:"created_at,notnull,default:current_timestamp"`
}
