package models

import "time"

type Task struct {
	ID        int       `json:"id" bun:"id,pk,autoincrement"`
	Name      string    `json:"name" bun:"name,notnull"`
	Status    string    `json:"status" bun:"status,notnull,default:'backlog'"`
	CreatedAt time.Time `json:"created_at" bun:"created_at,notnull,default:current_timestamp"`
	UpdatedAt time.Time `json:"updated_at" bun:"updated_at,notnull,default:current_timestamp"`
}
