package api

import (
	"net/http"

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
	handlers    *handlers.TaskHandlers
	middlewares *middlewareResource
}

func newTaskResource(operations handlers.TaskOperations, middlewares *middlewareResource) *taskResource {
	return &taskResource{
		handlers:    handlers.NewTaskHandlers(operations),
		middlewares: middlewares,
	}
}

type authResource struct {
	handlers *handlers.AuthHandlers
}

func newAuthResource(operations handlers.AuthOperations, userOperations handlers.UserOperations) *authResource {
	return &authResource{
		handlers: handlers.NewAuthHandlers(operations, userOperations),
	}
}

type middlewareResource struct {
	TokenClaimCtx func(http.Handler) http.Handler
	AdminAccess   func(http.Handler) http.Handler
}

func newMiddlewareResource(authStore *auth.AuthStore, userStore *database.UserStore) *middlewareResource {
	return &middlewareResource{
		TokenClaimCtx: newAuthResource(authStore, userStore).handlers.TokenClaimCtx,
		AdminAccess:   newAuthResource(authStore, userStore).handlers.AdminAccess,
	}
}

type API struct {
	Users *userResource
	Tasks *taskResource
	Auth  *authResource
}

func NewAPI(db *bun.DB, jwtAuth *jwtauth.JWTAuth) (*API, error) {
	authStore := auth.NewAuthStore(jwtAuth)
	userStore := database.NewUserStore(db)
	taskStore := database.NewTaskStore(db)

	middleware := newMiddlewareResource(authStore, userStore)
	auth := newAuthResource(authStore, userStore)
	users := newUserResource(userStore)
	tasks := newTaskResource(taskStore, middleware)

	api := &API{
		Users: users,
		Tasks: tasks,
		Auth:  auth,
	}

	return api, nil
}
