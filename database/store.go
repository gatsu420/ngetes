package database

import "github.com/uptrace/bun"

type UserStore struct {
	DB *bun.DB
}

func NewUserStore(db *bun.DB) *UserStore {
	return &UserStore{
		DB: db,
	}
}

type TaskStore struct {
	DB *bun.DB
}

func NewTaskStore(db *bun.DB) *TaskStore {
	return &TaskStore{
		DB: db,
	}
}

type UptimeStore struct {
	DB *bun.DB
}

func NewUptimeStore(db *bun.DB) *UptimeStore {
	return &UptimeStore{
		DB: db,
	}
}
