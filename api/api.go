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
	handlers      *handlers.TaskHandlers
	tokenClaimCtx func(http.Handler) http.Handler
	adminAccess   func(http.Handler) http.Handler
}

func newTaskResource(
	operations handlers.TaskOperations,
	tokenClaimCtx func(http.Handler) http.Handler,
	adminAccess func(http.Handler) http.Handler,
) *taskResource {
	return &taskResource{
		handlers:      handlers.NewTaskHandlers(operations),
		tokenClaimCtx: tokenClaimCtx,
		adminAccess:   adminAccess,
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

type API struct {
	Users *userResource
	Tasks *taskResource
	Auth  *authResource
}

func NewAPI(db *bun.DB, jwtAuth *jwtauth.JWTAuth) (*API, error) {
	authStore := auth.NewAuthStore(jwtAuth)
	userStore := database.NewUserStore(db)
	taskStore := database.NewTaskStore(db)

	auth := newAuthResource(authStore, userStore)
	users := newUserResource(userStore)
	tasks := newTaskResource(taskStore, auth.handlers.TokenClaimCtx, auth.handlers.AdminAccess)

	api := &API{
		Users: users,
		Tasks: tasks,
		Auth:  auth,
	}

	return api, nil
}
