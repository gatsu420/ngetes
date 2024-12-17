package main

import (
	"net/http"

	"github.com/gatsu420/ngetes/api"
	"github.com/gatsu420/ngetes/auth"
	"github.com/gatsu420/ngetes/database"
	"github.com/gatsu420/ngetes/logger"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
	"go.uber.org/zap"
)

func main() {
	err := logger.NewLogger()
	if err != nil {
		logger.Logger.Fatal("failed to initiate logger", zap.Error(err))
	}
	defer logger.Logger.Sync()

	db, err := database.DBConn()
	if err != nil {
		logger.Logger.Fatal("failed to connect to database", zap.Error(err))
	}
	defer db.Close()

	cache, err := database.RedisConn()
	if err != nil {
		logger.Logger.Fatal("failed to connect to redis", zap.Error(err))
	}
	defer cache.Close()

	auth, err := auth.JWTAuth()
	if err != nil {
		logger.Logger.Fatal("failed to generate JWT auth", zap.Error(err))
	}

	api, err := api.NewAPI(db, cache, auth)
	if err != nil {
		logger.Logger.Fatal("failed to initialize API", zap.Error(err))
	}

	go api.Uptime.Worker()

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
	logger.Logger.Info("starting server", zap.String("port", port))
	err = http.ListenAndServe(port, router)
	if err != nil {
		logger.Logger.Fatal("failed to start server", zap.Error(err))
	}
}
