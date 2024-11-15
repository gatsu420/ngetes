package database

import (
	"context"
	"net/url"
	"strconv"

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

func (s *TaskStore) List(f *TaskFilter) ([]models.Task, error) {
	t := []models.Task{}

	err := s.db.NewSelect().
		Model(&t).
		Apply(f.Apply).
		Scan(context.Background())
	if err != nil {
		return nil, err
	}

	return t, nil
}
