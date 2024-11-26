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

const ctxTask ctxKey = iota

type TaskOperations interface {
	List(*database.TaskFilter) ([]models.Task, error)
	Get(id int) (*models.Task, error)
	Create(*models.Task) (taskID int, err error)
	Update(*models.Task) error
	Delete(*models.Task) error

	CreateTracker(*models.Event) error
}

type TaskHandlers struct {
	Operations TaskOperations
}

func NewTaskHandlers(operations TaskOperations) *TaskHandlers {
	return &TaskHandlers{
		Operations: operations,
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

		task, err := rs.Operations.Get(id)
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

	task, err := rs.Operations.List(filters)
	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	event := &models.Event{
		TaskID: 0,
		Name:   "list",
	}
	err = rs.Operations.CreateTracker(event)
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
	err := rs.Operations.CreateTracker(event)
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

	taskID, err := rs.Operations.Create(task.Task)
	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	event := &models.Event{
		TaskID: taskID,
		Name:   "create",
	}
	err = rs.Operations.CreateTracker(event)
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

	err = rs.Operations.Update(task)
	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	event := &models.Event{
		TaskID: task.ID,
		Name:   "update",
	}
	err = rs.Operations.CreateTracker(event)
	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	render.Respond(w, r, NewTaskResponse(task))
}

func (rs *TaskHandlers) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	task := r.Context().Value(ctxTask).(*models.Task)

	err := rs.Operations.Delete(task)
	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	event := &models.Event{
		TaskID: task.ID,
		Name:   "delete",
	}
	err = rs.Operations.CreateTracker(event)
	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	render.Respond(w, r, NewTaskResponse(task))
}
