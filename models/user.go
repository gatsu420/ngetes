package models

import "time"

type User struct {
	ID        int       `json:"id" bun:"id,pk,autoincrement"`
	Name      string    `json:"name" bun:"name,notnull"`
	RoleID    int       `json:"-" bun:"role_id,notnull"`
	RoleName  string    `json:"role_name" bun:"-"`
	CreatedAt time.Time `json:"created_at" bun:"created_at,notnull,default:current_timestamp"`
	UpdatedAt time.Time `json:"updated_at" bun:"updated_at,notnull,default:current_timestamp"`
}

type Role struct {
	ID        int       `json:"id" bun:"id,pk,autoincrement"`
	Name      string    `json:"name" bun:"name,notnull"`
	CreatedAt time.Time `json:"created_at" bun:"created_at,notnull,default:current_timestamp"`
	UpdatedAt time.Time `json:"updated_at" bun:"updated_at,notnull,default:current_timestamp"`
}
