package api

import "github.com/go-chi/chi/v5"

func (rs *userResource) Router() *chi.Mux {
	router := chi.NewRouter()

	router.Post("/", rs.handlers.CreateUserHandler)

	return router
}
