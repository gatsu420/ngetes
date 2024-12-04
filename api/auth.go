package api

import "github.com/go-chi/chi/v5"

func (rs *authResource) LandingRouter() *chi.Mux {
	router := chi.NewRouter()

	router.Get("/", rs.handlers.GetTokenHandler)

	return router
}

func (rs *authResource) Router() *chi.Mux {
	router := chi.NewRouter()

	router.Group(func(router chi.Router) {
		router.Use(rs.handlers.TokenClaimCtx)
		router.Get("/claim", rs.handlers.GetTokenClaimHandler)
	})

	return router
}
