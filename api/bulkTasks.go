package api

import "github.com/go-chi/chi/v5"

func (r *bulkTasksResource) Router() *chi.Mux {
	router := chi.NewRouter()

	router.Use(r.middlewares.TokenClaimCtx)
	router.Use(r.middlewares.TokenBlacklistAccess)
	router.Post("/", r.handlers.CreateHandler)

	return router
}
