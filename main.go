package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gatsu420/ngetes/internal/api"
	"github.com/gatsu420/ngetes/internal/auth"
	"github.com/gatsu420/ngetes/internal/database"
	"github.com/gatsu420/ngetes/internal/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
)

func main() {
	db, err := database.DBConn()
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	api, err := api.NewAPI(db)
	if err != nil {
		log.Fatalf("failed to initialize API: %v", err)
	}

	tokenAuth := auth.NewAuth()

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Heartbeat("/ping"))

	// TODO: Refactor handlers
	router.Post("/login", handlers.LoginHandler)
	router.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(jwtauth.Authenticator(tokenAuth))
		r.Mount("/tasks", api.Tasks.Router())
		r.Get("/hello", func(w http.ResponseWriter, r *http.Request) {
			_, claims, _ := jwtauth.FromContext(r.Context())

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"claims": claims,
			})
		})
	})

	port := ":8080"
	fmt.Println("starting server on port", port)
	err = http.ListenAndServe(port, router)
	if err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
