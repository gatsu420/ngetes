package api

import (
	"github.com/gatsu420/ngetes/internal/database"
	"github.com/gatsu420/ngetes/internal/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/uptrace/bun"
)

// type ctxKey int

// const (
// 	ctxTask ctxKey = iota
// )

// type TaskStore interface {
// 	List(*database.TaskFilter) ([]models.Task, error)
// 	Get(id int) (*models.Task, error)
// 	Create(*models.Task) (taskID int, err error)
// 	Update(*models.Task) error
// 	Delete(*models.Task) error

// 	CreateTracker(*models.Event) error
// }

type TaskResource struct {
	Handlers *handlers.TaskHandlers
}

func NewTaskResource(s handlers.TaskStore) *TaskResource {
	return &TaskResource{
		Handlers: handlers.NewTaskHandler(s),
	}
}

type API struct {
	Tasks *TaskResource
}

func NewAPI(db *bun.DB) (*API, error) {
	taskStore := database.NewTaskStore(db)
	tasks := NewTaskResource(taskStore)
	api := &API{
		Tasks: tasks,
	}

	return api, nil
}

func (rs *TaskResource) Router() *chi.Mux {
	router := chi.NewRouter()

	router.Get("/", rs.Handlers.ListHandler)
	router.Post("/", rs.Handlers.CreateHandler)
	router.Route("/{taskID}", func(router chi.Router) {
		router.Use(rs.Handlers.TaskCtx)
		router.Get("/", rs.Handlers.GetHandler)
		router.Put("/", rs.Handlers.UpdateHandler)
		router.Delete("/", rs.Handlers.DeleteHandler)
	})

	return router
}
