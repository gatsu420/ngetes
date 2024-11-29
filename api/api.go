package api

import (
	"github.com/gatsu420/ngetes/auth"
	"github.com/gatsu420/ngetes/database"
	"github.com/gatsu420/ngetes/handlers"
	"github.com/go-chi/jwtauth/v5"
	"github.com/uptrace/bun"
)

type userResource struct {
	handlers *handlers.UserHandlers
}

func newUserResource(operations handlers.UserOperations) *userResource {
	return &userResource{
		handlers: handlers.NewUserHandlers(operations),
	}
}

type taskResource struct {
	handlers *handlers.TaskHandlers
}

func newTaskResource(operations handlers.TaskOperations) *taskResource {
	return &taskResource{
		handlers: handlers.NewTaskHandlers(operations),
	}
}

type authResource struct {
	handlers *handlers.AuthHandlers
}

func newAuthResource(operations handlers.AuthOperations) *authResource {
	return &authResource{
		handlers: handlers.NewAuthHandlers(operations),
	}
}

type API struct {
	Users *userResource
	Tasks *taskResource
	Auth  *authResource
}

func NewAPI(db *bun.DB, jwtAuth *jwtauth.JWTAuth) (*API, error) {
	userStore := database.NewUserStore(db)
	taskStore := database.NewTaskStore(db)
	authStore := auth.NewAuthStore(jwtAuth)

	users := newUserResource(userStore)
	tasks := newTaskResource(taskStore)
	auth := newAuthResource(authStore)

	api := &API{
		Users: users,
		Tasks: tasks,
		Auth:  auth,
	}

	return api, nil
}
