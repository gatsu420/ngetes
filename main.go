package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gatsu420/ngetes/api"
	"github.com/gatsu420/ngetes/auth"
	"github.com/gatsu420/ngetes/database"
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

	rdb, err := database.RedisConn()
	if err != nil {
		log.Fatalf("failed to connect to redis: %v", err)
	}
	defer rdb.Close()

	auth, err := auth.JWTAuth()
	if err != nil {
		log.Fatalf("failed to generate JWT auth: %v", err)
	}

	api, err := api.NewAPI(db, auth)
	if err != nil {
		log.Fatalf("failed to initialize API: %v", err)
	}

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Heartbeat("/ping"))

	router.Mount("/users", api.Users.Router())
	router.Mount("/token", api.Auth.LandingRouter())
	router.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(auth))
		r.Use(jwtauth.Authenticator(auth))
		r.Mount("/auth", api.Auth.Router())
		r.Mount("/tasks", api.Tasks.Router())
	})

	port := ":8080"
	fmt.Println("starting server on port", port)
	err = http.ListenAndServe(port, router)
	if err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
