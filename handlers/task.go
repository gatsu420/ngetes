package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gatsu420/ngetes/database"
	"github.com/gatsu420/ngetes/models"
	"github.com/go-chi/render"
)

type TaskOperations interface {
	List(*database.TaskFilter) ([]models.Task, error)
	Get(id int) (*models.Task, error)
	Create(*models.Task) (taskID int, err error)
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

func GetWeatherForecast(url string, fc <-chan bool) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	time.Sleep(5 * time.Second)
	fmt.Println(string(body))

	if ok := <-fc; ok {
		fmt.Println("table doesnt exist")
	}

	return nil
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

	failedTrackerChan := make(chan bool)
	go func() {
		event := &models.Event{
			TaskID: 0,
			Name:   "list",
		}
		err := hd.Operations.CreateTracker(event)
		if err != nil {
			log.Printf("failed to create tracker event: %v", err)
			failedTrackerChan <- true
		}
	}()

	forecast_url := "https://api.open-meteo.com/v1/forecast?latitude=-6.4&longitude=106.8186&hourly=temperature_2m"
	go GetWeatherForecast(forecast_url, failedTrackerChan)

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
			log.Printf("failed to create tracker event: %v", err)
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
			log.Printf("failed to create tracker event: %v", err)
		}
	}()

	render.Respond(w, r, newTaskResponse(task.Task))
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
			log.Printf("failed to create tracker event: %v", err)
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
			log.Printf("failed to create tracker event: %v", err)
		}
	}()

	render.Respond(w, r, newDeletedTaskResponse(task))
}
