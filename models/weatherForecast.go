package models

import (
	"time"

	"github.com/uptrace/bun"
)

type WeatherForecast struct {
	bun.BaseModel `bun:"table:weather_forecast"`

	ID         int       `json:"id" bun:"id,pk,autoincrement"`
	StatusCode int       `json:"status_code" bun:"status_code,notnull"`
	URL        string    `json:"url" bun:"url,notnull"`
	Payload    string    `json:"payload" bun:"payload,notnull"`
	CreatedAt  time.Time `json:"created_at" bun:"created_at,notnull,default:current_timestamp"`
}
