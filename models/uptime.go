package models

import (
	"time"

	"github.com/uptrace/bun"
)

type Uptime struct {
	bun.BaseModel `bun:"table:uptime"`

	ID        int       `json:"id" bun:"id,pk,autoincrement"`
	CreatedAt time.Time `json:"created_at" bun:"created_at,notnull,default:current_timestamp"`
}
