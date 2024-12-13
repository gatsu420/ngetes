package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/gatsu420/ngetes/models"
)

func (s *WeatherForecastStore) CreateForecast(f *models.WeatherForecast) (
	id int, createdAt time.Time, err error,
) {
	ctx := context.Background()
	tx, err := s.DB.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return 0, time.Time{}, err
	}

	_, err = tx.NewInsert().
		Model(f).
		Exec(ctx)
	if err != nil {
		tx.Rollback()
		return 0, time.Time{}, err
	}

	tx.Commit()
	return f.ID, f.CreatedAt, nil
}
