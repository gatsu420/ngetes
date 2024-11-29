package api

import "github.com/go-chi/chi/v5"

func (rs *authResource) Router() *chi.Mux {
	router := chi.NewRouter()

	router.Get("/token", rs.handlers.LoginHandler)

	return router
}
