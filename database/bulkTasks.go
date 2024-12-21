package database

import (
	"context"
	"database/sql"

	"github.com/gatsu420/ngetes/models"
)

func (s *BulkTasksStore) Create(t []models.Task) (tasks []models.Task, err error) {
	ctx := context.Background()
	tx, err := s.DB.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil, err
	}

	_, err = tx.NewInsert().
		Model(&t).
		Exec(ctx)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	for _, task := range t {
		tasks = append(tasks, task)
	}

	tx.Commit()
	return tasks, nil
}
