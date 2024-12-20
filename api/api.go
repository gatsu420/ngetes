package api

import (
	"net/http"

	"github.com/gatsu420/ngetes/auth"
	"github.com/gatsu420/ngetes/database"
	"github.com/gatsu420/ngetes/handlers"
	"github.com/gatsu420/ngetes/workers"
	"github.com/go-chi/jwtauth/v5"
	"github.com/redis/go-redis/v9"
	"github.com/uptrace/bun"
)

type authResource struct {
	handlers *handlers.AuthHandlers
}

func newAuthResource(operations handlers.AuthOperations, userOperations handlers.UserOperations) *authResource {
	return &authResource{
		handlers: handlers.NewAuthHandlers(operations, userOperations),
	}
}

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

type bulkTasksResource struct {
	handlers    *handlers.BulkTasksHandlers
	middlewares *middlewareResource
}

func newBulkTasksResource(operations handlers.BulkTasksOperations, middlewares *middlewareResource) *bulkTasksResource {
	return &bulkTasksResource{
		handlers:    handlers.NewBulkTasksHandlers(operations),
		middlewares: middlewares,
	}
}

type middlewareResource struct {
	TokenClaimCtx        func(http.Handler) http.Handler
	TokenBlacklistAccess func(http.Handler) http.Handler
	AdminAccess          func(http.Handler) http.Handler
}

func newMiddlewareResource(authStore *auth.AuthStore, userStore *database.UserStore) *middlewareResource {
	return &middlewareResource{
		TokenClaimCtx:        newAuthResource(authStore, userStore).handlers.TokenClaimCtx,
		TokenBlacklistAccess: newAuthResource(authStore, userStore).handlers.TokenBlacklistAccess,
		AdminAccess:          newAuthResource(authStore, userStore).handlers.AdminAccess,
	}
}

type uptimeResource struct {
	workers *workers.UptimeWorkers
}

func newUptimeResource(operations workers.UptimeOperations) *uptimeResource {
	return &uptimeResource{
		workers: workers.NewUptimeWorkers(operations),
	}
}

type API struct {
	Users     *userResource
	Tasks     *taskResource
	BulkTasks *bulkTasksResource
	Auth      *authResource
	Uptime    *uptimeResource
}

func NewAPI(db *bun.DB, cache *redis.Client, jwtAuth *jwtauth.JWTAuth) (*API, error) {
	authStore := auth.NewAuthStore(jwtAuth, cache)
	userStore := database.NewUserStore(db)
	taskStore := database.NewTaskStore(db)
	bulkTasksStore := database.NewBulkTasksStore(db)
	uptimeStore := database.NewUptimeStore(db)

	middleware := newMiddlewareResource(authStore, userStore)
	auth := newAuthResource(authStore, userStore)
	users := newUserResource(userStore)
	tasks := newTaskResource(taskStore, middleware)
	bulkTasks := newBulkTasksResource(bulkTasksStore, middleware)
	uptime := newUptimeResource(uptimeStore)

	api := &API{
		Users:     users,
		Tasks:     tasks,
		BulkTasks: bulkTasks,
		Auth:      auth,
		Uptime:    uptime,
	}

	return api, nil
}
