package api

import "github.com/go-chi/chi/v5"

func (rs *taskResource) Router() *chi.Mux {
	router := chi.NewRouter()

	router.Use(rs.middlewares.TokenClaimCtx)
	router.Use(rs.middlewares.TokenBlacklistAccess)
	router.Get("/", rs.handlers.ListHandler)
	router.Post("/", rs.handlers.CreateHandler)
	router.Post("/bulk", rs.handlers.CreateBulkHandler)
	router.Route("/{taskID}", func(router chi.Router) {
		router.Use(rs.handlers.TaskCtx)
		router.Get("/", rs.handlers.GetHandler)

		router.Group(func(router chi.Router) {
			router.Use(rs.middlewares.AdminAccess)
			router.Put("/", rs.handlers.UpdateHandler)
			router.Delete("/", rs.handlers.DeleteHandler)
		})
	})

	return router
}
