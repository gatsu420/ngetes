package handlers

import (
	"net/http"
	"time"

	"github.com/gatsu420/ngetes/database"
	"github.com/gatsu420/ngetes/logger"
	"github.com/gatsu420/ngetes/models"
	"github.com/go-chi/render"
	"github.com/xuri/excelize/v2"
	"go.uber.org/zap"
)

type TaskOperations interface {
	List(*database.TaskFilter) ([]models.Task, error)
	Get(id int) (*models.Task, error)
	Create(*models.Task) (taskID int, err error)
	CreateBulk(t []models.Task) (tasks []models.Task, err error)
	Update(*models.Task) error
	Delete(*models.Task) error

	CreateTracker(*models.Event) error

	GetClaim(r *http.Request) (map[string]interface{}, error)
}

type TaskHandlers struct {
	Operations TaskOperations
}

func NewTaskHandlers(operations TaskOperations) *TaskHandlers {
	return &TaskHandlers{
		Operations: operations,
	}
}

func (hd *TaskHandlers) ListHandler(w http.ResponseWriter, r *http.Request) {
	filters, err := database.NewTaskFilter(r.URL.Query())
	if err != nil {
		render.Render(w, r, errRender(err))
		return
	}

	task, err := hd.Operations.List(filters)
	if err != nil {
		render.Render(w, r, errRender(err))
		return
	}

	go func() {
		event := &models.Event{
			TaskID: 0,
			Name:   "list",
		}
		err := hd.Operations.CreateTracker(event)
		if err != nil {
			logger.Logger.Error("failed to create tracker event", zap.Error(err))
		}
	}()

	render.Respond(w, r, newTaskListResponse(&task))
}

func (hd *TaskHandlers) GetHandler(w http.ResponseWriter, r *http.Request) {
	task := r.Context().Value(taskCtx).(*models.Task)

	go func() {
		event := &models.Event{
			TaskID: task.ID,
			Name:   "get",
		}
		err := hd.Operations.CreateTracker(event)
		if err != nil {
			logger.Logger.Error("failed to create tracker event", zap.Error(err))
		}
	}()

	render.Respond(w, r, newTaskResponse(task))
}

func (hd *TaskHandlers) CreateHandler(w http.ResponseWriter, r *http.Request) {
	task := &taskRequest{}
	err := render.Bind(r, task)
	if err != nil {
		render.Render(w, r, errRender(err))
		return
	}

	taskID, err := hd.Operations.Create(task.Task)
	if err != nil {
		render.Render(w, r, errRender(err))
		return
	}

	go func() {
		event := &models.Event{
			TaskID: taskID,
			Name:   "create",
		}
		err = hd.Operations.CreateTracker(event)
		if err != nil {
			logger.Logger.Error("failed to create tracker event", zap.Error(err))
		}
	}()

	render.Respond(w, r, newTaskResponse(task.Task))
}

func (hd *TaskHandlers) CreateBulkHandler(w http.ResponseWriter, r *http.Request) {
	f, err := excelize.OpenFile("Book1.xlsx")
	if err != nil {
		render.Render(w, r, errRender(err))
		return
	}
	defer f.Close()

	rows, err := f.GetRows("Sheet1")
	if err != nil {
		render.Render(w, r, errRender(err))
		return
	}

	tasks := []models.Task{}
	for i, row := range rows {
		if i == 0 {
			continue
		}

		task := models.Task{
			Name:   row[0],
			Status: row[1],
		}
		tasks = append(tasks, task)
	}

	t, err := hd.Operations.CreateBulk(tasks)
	if err != nil {
		render.Render(w, r, errRender(err))
		return
	}

	render.Respond(w, r, newTaskListResponse(&t))
}

func (hd *TaskHandlers) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	task := r.Context().Value(taskCtx).(*models.Task)
	task.UpdatedAt = time.Now()

	data := &taskRequest{
		Task: task,
	}
	err := render.Bind(r, data)
	if err != nil {
		render.Render(w, r, errRender(err))
		return
	}

	err = hd.Operations.Update(task)
	if err != nil {
		render.Render(w, r, errRender(err))
		return
	}

	go func() {
		event := &models.Event{
			TaskID: task.ID,
			Name:   "update",
		}
		err = hd.Operations.CreateTracker(event)
		if err != nil {
			logger.Logger.Error("failed to create tracker event", zap.Error(err))
		}
	}()

	render.Respond(w, r, newTaskResponse(task))
}

func (hd *TaskHandlers) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	task := r.Context().Value(taskCtx).(*models.Task)

	err := hd.Operations.Delete(task)
	if err != nil {
		render.Render(w, r, errRender(err))
		return
	}

	go func() {
		event := &models.Event{
			TaskID: task.ID,
			Name:   "delete",
		}
		err = hd.Operations.CreateTracker(event)
		if err != nil {
			logger.Logger.Error("failed to create tracker event", zap.Error(err))
		}
	}()

	render.Respond(w, r, newDeletedTaskResponse(task))
}
