package database

import "github.com/uptrace/bun"

type taskStore struct {
	DB *bun.DB
}

func NewTaskStore(db *bun.DB) *taskStore {
	return &taskStore{
		DB: db,
	}
}
