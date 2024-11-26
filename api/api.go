package api

import (
	"github.com/gatsu420/ngetes/auth"
	"github.com/gatsu420/ngetes/database"
	"github.com/gatsu420/ngetes/handlers"
	"github.com/go-chi/jwtauth/v5"
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

type loginResource struct {
	Handlers *handlers.LoginHandlers
}

func NewLoginResource(operations handlers.LoginOperations) *loginResource {
	return &loginResource{
		Handlers: handlers.NewLoginHandler(operations),
	}
}

type API struct {
	Tasks *taskResource
	Login *loginResource
}

func NewAPI(db *bun.DB, jwtAuth *jwtauth.JWTAuth) (*API, error) {
	taskStore := database.NewTaskStore(db)
	authStore := auth.NewAuthStore(jwtAuth)

	tasks := NewTaskResource(taskStore)
	auth := NewLoginResource(authStore)

	api := &API{
		Tasks: tasks,
		Login: auth,
	}

	return api, nil
}
