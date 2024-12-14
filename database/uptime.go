package database

import (
	"context"
	"database/sql"

	"github.com/gatsu420/ngetes/models"
)

func (s *UptimeStore) CreateUptime(u *models.Uptime) error {
	ctx := context.Background()
	tx, err := s.DB.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}

	_, err = tx.NewInsert().
		Model(u).
		Exec(ctx)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}
