package database

import "github.com/uptrace/bun"

type userStore struct {
	DB *bun.DB
}

func NewUserStore(db *bun.DB) *userStore {
	return &userStore{
		DB: db,
	}
}

type taskStore struct {
	DB *bun.DB
}

func NewTaskStore(db *bun.DB) *taskStore {
	return &taskStore{
		DB: db,
	}
}
