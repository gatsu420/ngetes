package database

import (
	"context"
	"database/sql"

	"github.com/gatsu420/ngetes/models"
)

func (s *UserStore) CreateUser(u *models.User) error {
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

func (s *UserStore) ListRoles() ([]models.Role, error) {
	ctx := context.Background()
	tx, err := s.DB.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil, err
	}

	roles := []models.Role{}
	err = s.DB.NewSelect().
		Model(&roles).
		Distinct().
		Column("name").
		Scan(ctx)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	return roles, nil
}

func (s *UserStore) GetRoleID(roleModel *models.Role, roleName string) (roleID int, err error) {
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

func (s *UserStore) GetUserRole(name string) (roleID int, err error) {
	ctx := context.Background()
	tx, err := s.DB.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return 0, err
	}

	err = s.DB.NewSelect().
		Model((*models.User)(nil)).
		Column("role_id").
		Where("name = ?", name).
		Scan(ctx, &roleID)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	tx.Commit()
	return roleID, nil
}

func (s *UserStore) GetValidUserName(userName string) (isExist bool, err error) {
	ctx := context.Background()
	tx, err := s.DB.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return false, err
	}

	existence, err := s.DB.NewSelect().
		Model((*models.User)(nil)).
		Distinct().
		Column("name").
		Where("name = ?", userName).
		Exists(ctx)
	if err != nil {
		tx.Rollback()
		return false, err
	}

	tx.Commit()
	return existence, nil
}
