package database

import (
	"context"
	"database/sql"
	"net/http"
	"net/url"
	"strconv"

	"github.com/gatsu420/ngetes/models"
	"github.com/go-chi/jwtauth/v5"
	"github.com/uptrace/bun"
)

type TaskFilter struct {
	Limit  int
	Offset int
	Filter map[string]interface{}
	Order  []string
}

func NewTaskFilter(params interface{}) (*TaskFilter, error) {
	v, _ := params.(url.Values)
	f := &TaskFilter{
		Limit:  2,
		Offset: 1,
		Filter: map[string]interface{}{},
		Order:  v["order"],
	}

	if limit := v.Get("limit"); limit != "" {
		f.Limit, _ = strconv.Atoi(limit)
	}

	if offset := v.Get("offset"); offset != "" {
		f.Offset, _ = strconv.Atoi(offset)
	}

	for key, value := range v {
		if key != "limit" && key != "offset" && key != "order" {
			f.Filter[key] = value[0]
		}
	}

	return f, nil
}

func (f *TaskFilter) Apply(q *bun.SelectQuery) *bun.SelectQuery {
	q = q.Limit(f.Limit).Offset(f.Offset)

	for k, v := range f.Filter {
		q = q.Where("? = ?", bun.Ident(k), v)
	}

	for _, v := range f.Order {
		q = q.Order(v)
	}

	return q
}

func (s *taskStore) List(f *TaskFilter) ([]models.Task, error) {
	t := []models.Task{}

	err := s.DB.NewSelect().
		Model(&t).
		Apply(f.Apply).
		Scan(context.Background())
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (s *taskStore) Get(id int) (*models.Task, error) {
	ctx := context.Background()
	tx, err := s.DB.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil, err
	}

	t := &models.Task{
		ID: id,
	}
	err = s.DB.NewSelect().
		Model(t).
		WherePK().
		Scan(ctx)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()
	return t, nil
}

func (s *taskStore) Create(t *models.Task) (taskID int, err error) {
	ctx := context.Background()
	tx, err := s.DB.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return 0, err
	}

	_, err = tx.NewInsert().
		Model(t).
		Exec(ctx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	tx.Commit()
	return t.ID, nil
}

func (s *taskStore) Update(t *models.Task) error {
	ctx := context.Background()
	tx, err := s.DB.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}

	_, err = s.DB.NewUpdate().
		Model(t).
		WherePK().
		Exec(ctx)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (s *taskStore) Delete(t *models.Task) error {
	ctx := context.Background()
	tx, err := s.DB.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}

	_, err = s.DB.NewDelete().
		Model(t).
		WherePK().
		Exec(ctx)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (s *taskStore) CreateTracker(e *models.Event) error {
	ctx := context.Background()
	tx, err := s.DB.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}

	event := &models.Event{
		TaskID: e.TaskID,
		Name:   e.Name,
	}
	_, err = s.DB.NewInsert().
		Model(event).
		Exec(ctx)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (s *taskStore) GetClaim(r *http.Request) (map[string]interface{}, error) {
	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil {
		return nil, err
	}

	return claims, nil
}
