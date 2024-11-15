package api

import (
	"net/http"

	"github.com/gatsu420/ngetes/internal/database"
	"github.com/gatsu420/ngetes/internal/handlers"
	"github.com/gatsu420/ngetes/models"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/uptrace/bun"
)

type TaskStore interface {
	List(*database.TaskFilter) ([]models.Task, error)
}

type TaskResource struct {
	Store TaskStore
}

func NewTaskResource(s TaskStore) *TaskResource {
	return &TaskResource{
		Store: s,
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

	router.Get("/", rs.list)

	return router
}

func (rs *TaskResource) list(w http.ResponseWriter, r *http.Request) {
	f, _ := database.NewTaskFilter(r.URL.Query())
	acc, err := rs.Store.List(f)
	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	render.Respond(w, r, handlers.NewTaskListResponse(&acc))
}
