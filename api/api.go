package api

import (
	"github.com/gatsu420/ngetes/database"
	"github.com/gatsu420/ngetes/handlers"
	"github.com/uptrace/bun"
)

type taskResource struct {
	Handlers *handlers.TaskHandlers
}

func NewTaskResource(operations handlers.TaskOperations) *taskResource {
	return &taskResource{
		Handlers: handlers.NewTaskHandler(operations),
	}
}

type API struct {
	Tasks *taskResource
}

func NewAPI(db *bun.DB) (*API, error) {
	taskStore := database.NewTaskStore(db)
	tasks := NewTaskResource(taskStore)

	api := &API{
		Tasks: tasks,
	}

	return api, nil
}
