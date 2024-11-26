package api

import "github.com/go-chi/chi/v5"

func (rs *taskResource) Router() *chi.Mux {
	router := chi.NewRouter()

	router.Get("/", rs.Handlers.ListHandler)
	router.Post("/", rs.Handlers.CreateHandler)
	router.Route("/{taskID}", func(router chi.Router) {
		router.Use(rs.Handlers.TaskCtx)
		router.Get("/", rs.Handlers.GetHandler)
		router.Put("/", rs.Handlers.UpdateHandler)
		router.Delete("/", rs.Handlers.DeleteHandler)
	})

	router.Get("/claims", rs.Handlers.GetClaimHandler)

	return router
}
