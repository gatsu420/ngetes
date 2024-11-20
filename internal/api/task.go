package api

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gatsu420/ngetes/internal/database"
	"github.com/gatsu420/ngetes/internal/handlers"
	"github.com/gatsu420/ngetes/models"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/uptrace/bun"
)

type ctxKey int

const (
	ctxTask ctxKey = iota
)

type TaskStore interface {
	List(*database.TaskFilter) ([]models.Task, error)
	Get(id int) (*models.Task, error)
	Create(*models.Task) (taskID int, err error)
	Update(*models.Task) error
	Delete(*models.Task) error

	CreateTracker(*models.Event) error
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
	router.Post("/", rs.create)
	router.Route("/{taskID}", func(router chi.Router) {
		router.Use(rs.taskCtx)
		router.Get("/", rs.get)
		router.Put("/", rs.update)
		router.Delete("/", rs.delete)
	})

	return router
}

func (rs *TaskResource) taskCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "taskID"))
		if err != nil {
			render.Render(w, r, ErrRender(err))
			return
		}
		acc, err := rs.Store.Get(id)
		if err != nil {
			render.Render(w, r, ErrRender(err))
			return
		}

		ctx := context.WithValue(r.Context(), ctxTask, acc)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
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

func (rs *TaskResource) get(w http.ResponseWriter, r *http.Request) {
	task := r.Context().Value(ctxTask).(*models.Task)

	event := &models.Event{
		TaskID: task.ID,
		Name:   "Get",
	}
	err := rs.Store.CreateTracker(event)
	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	render.Respond(w, r, handlers.NewTaskResponse(task))
}

func (rs *TaskResource) create(w http.ResponseWriter, r *http.Request) {
	task := &handlers.TaskRequest{}
	err := render.Bind(r, task)
	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	taskID, err := rs.Store.Create(task.Task)
	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	event := &models.Event{
		TaskID: taskID,
		Name:   "Create",
	}
	err = rs.Store.CreateTracker(event)
	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	render.Respond(w, r, handlers.NewTaskResponse(task.Task))
}

func (rs *TaskResource) update(w http.ResponseWriter, r *http.Request) {
	task := r.Context().Value(ctxTask).(*models.Task)
	data := &handlers.TaskRequest{
		Task: task,
	}
	err := render.Bind(r, data)
	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	err = rs.Store.Update(task)
	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	event := &models.Event{
		TaskID: task.ID,
		Name:   "Update",
	}
	err = rs.Store.CreateTracker(event)
	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	render.Respond(w, r, handlers.NewTaskResponse(task))
}

func (rs *TaskResource) delete(w http.ResponseWriter, r *http.Request) {
	task := r.Context().Value(ctxTask).(*models.Task)

	err := rs.Store.Delete(task)
	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	render.Respond(w, r, handlers.NewTaskResponse(task))
}
