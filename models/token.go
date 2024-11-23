package models

import "time"

type Token struct {
	ID        int       `json:"id" bun:"id,pk,autoincrement"`
	Token     string    `json:"token" bun:"token,notnull"`
	ExpiredAt time.Time `json:"expired_at" bun:"expired_at,notnull"`
	CreatedAt time.Time `json:"created_at" bun:"created_at,notnull,default:current_timestamp"`
	UpdatedAt time.Time `json:"updated_at" bun:"updated_at,notnull,default:current_timestamp"`
}
