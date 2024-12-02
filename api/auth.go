package api

import "github.com/go-chi/chi/v5"

func (rs *authResource) Router() *chi.Mux {
	router := chi.NewRouter()

	router.Get("/token", rs.handlers.GetTokenHandler)
	router.Post("/login", rs.handlers.GetUserNameExistenceHandler)

	return router
}
