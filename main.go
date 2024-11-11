package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gatsu420/ngetes/internal/api"
	"github.com/gatsu420/ngetes/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Heartbeat("/ping"))
	router.Mount("/tasks", api.Tasks.Router())

	port := ":8080"
	fmt.Println("starting server on port", port)
	err = http.ListenAndServe(port, router)
	if err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
