package database

import (
	"context"
	"database/sql"

	"github.com/gatsu420/ngetes/models"
)

func (s *userStore) CreateUser(u *models.User) error {
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

func (s *userStore) GetUserRole(roleModel *models.Role, roleName string) (roleID int, err error) {
	ctx := context.Background()
	tx, err := s.DB.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return 0, err
	}

	err = s.DB.NewSelect().
		Model(roleModel).
		Where("name = ?", roleName).
		Scan(ctx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	tx.Commit()
	return roleModel.ID, nil
}
