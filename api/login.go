package api

import "github.com/go-chi/chi/v5"

func (rs *loginResource) Router() *chi.Mux {
	router := chi.NewRouter()

	router.Get("/", rs.handlers.LoginHandler)

	return router
}
