package database

import (
	"context"

	"github.com/gatsu420/ngetes/models"
	"github.com/uptrace/bun"
)

type TaskStore struct {
	db *bun.DB
}

func NewTaskStore(db *bun.DB) *TaskStore {
	return &TaskStore{
		db: db,
	}
}

func (s *TaskStore) List() ([]models.Task, error) {
	t := []models.Task{}

	err := s.db.NewSelect().
		Model(&t).
		Scan(context.Background())
	if err != nil {
		return nil, err
	}

	return t, nil
}
