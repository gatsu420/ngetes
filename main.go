package main

import (
	"net/http"

	"github.com/gatsu420/ngetes/api"
	"github.com/gatsu420/ngetes/auth"
	"github.com/gatsu420/ngetes/config"
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

	config, err := config.LoadConfig()
	if err != nil {
		logger.Logger.Fatal("failed to read config file", zap.Error(err))
	}

	db, err := database.DBConn(config)
	if err != nil {
		logger.Logger.Fatal("failed to connect to database", zap.Error(err))
	}
	defer db.Close()

	cache, err := database.RedisConn(config)
	if err != nil {
		logger.Logger.Fatal("failed to connect to redis", zap.Error(err))
	}
	defer cache.Close()

	auth, err := auth.JWTAuth(config)
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
		r.Mount("/bulktasks", api.BulkTasks.Router())
	})

	port := ":8080"
	logger.Logger.Info("starting server", zap.String("port", port))
	err = http.ListenAndServe(port, router)
	if err != nil {
		logger.Logger.Fatal("failed to start server", zap.Error(err))
	}
}
