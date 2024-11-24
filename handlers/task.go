package handlers

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gatsu420/ngetes/database"
	"github.com/gatsu420/ngetes/models"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
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

type TaskHandlers struct {
	Store TaskStore
}

func NewTaskHandler(s TaskStore) *TaskHandlers {
	return &TaskHandlers{
		Store: s,
	}
}

type taskListResponse struct {
	Task *[]models.Task `json:"tasks"`
}

func NewTaskListResponse(t *[]models.Task) *taskListResponse {
	return &taskListResponse{
		Task: t,
	}
}

type taskResponse struct {
	Task *models.Task `json:"task"`
}

func NewTaskResponse(t *models.Task) *taskResponse {
	return &taskResponse{
		Task: t,
	}
}

type TaskRequest struct {
	Task *models.Task `json:"task"`
}

func (tr *TaskRequest) Bind(r *http.Request) error {
	return nil
}

func (rs *TaskHandlers) TaskCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "taskID"))
		if err != nil {
			render.Render(w, r, ErrRender(err))
			return
		}

		task, err := rs.Store.Get(id)
		if err != nil {
			render.Render(w, r, ErrRender(err))
			return
		}

		ctx := context.WithValue(r.Context(), ctxTask, task)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (rs *TaskHandlers) ListHandler(w http.ResponseWriter, r *http.Request) {
	filters, err := database.NewTaskFilter(r.URL.Query())
	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	task, err := rs.Store.List(filters)
	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	event := &models.Event{
		TaskID: 0,
		Name:   "list",
	}
	err = rs.Store.CreateTracker(event)
	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	render.Respond(w, r, NewTaskListResponse(&task))
}

func (rs *TaskHandlers) GetHandler(w http.ResponseWriter, r *http.Request) {
	task := r.Context().Value(ctxTask).(*models.Task)

	event := &models.Event{
		TaskID: task.ID,
		Name:   "get",
	}
	err := rs.Store.CreateTracker(event)
	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	render.Respond(w, r, NewTaskResponse(task))
}

func (rs *TaskHandlers) CreateHandler(w http.ResponseWriter, r *http.Request) {
	task := &TaskRequest{}
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
		Name:   "create",
	}
	err = rs.Store.CreateTracker(event)
	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	render.Respond(w, r, NewTaskResponse(task.Task))
}

func (rs *TaskHandlers) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	task := r.Context().Value(ctxTask).(*models.Task)
	task.UpdatedAt = time.Now()

	data := &TaskRequest{
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
		Name:   "update",
	}
	err = rs.Store.CreateTracker(event)
	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	render.Respond(w, r, NewTaskResponse(task))
}

func (rs *TaskHandlers) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	task := r.Context().Value(ctxTask).(*models.Task)

	err := rs.Store.Delete(task)
	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	event := &models.Event{
		TaskID: task.ID,
		Name:   "delete",
	}
	err = rs.Store.CreateTracker(event)
	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	render.Respond(w, r, NewTaskResponse(task))
}
