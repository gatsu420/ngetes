package api

import (
	"github.com/gatsu420/ngetes/auth"
	"github.com/gatsu420/ngetes/database"
	"github.com/gatsu420/ngetes/handlers"
	"github.com/go-chi/jwtauth/v5"
	"github.com/uptrace/bun"
)

type taskResource struct {
	handlers *handlers.TaskHandlers
}

func newTaskResource(operations handlers.TaskOperations) *taskResource {
	return &taskResource{
		handlers: handlers.NewTaskHandlers(operations),
	}
}

type loginResource struct {
	handlers *handlers.LoginHandlers
}

func newLoginResource(operations handlers.LoginOperations) *loginResource {
	return &loginResource{
		handlers: handlers.NewLoginHandlers(operations),
	}
}

type API struct {
	Tasks *taskResource
	Login *loginResource
}

func NewAPI(db *bun.DB, jwtAuth *jwtauth.JWTAuth) (*API, error) {
	taskStore := database.NewTaskStore(db)
	authStore := auth.NewAuthStore(jwtAuth)

	tasks := newTaskResource(taskStore)
	auth := newLoginResource(authStore)

	api := &API{
		Tasks: tasks,
		Login: auth,
	}

	return api, nil
}
