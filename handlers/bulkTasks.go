package handlers

import (
	"net/http"

	"github.com/gatsu420/ngetes/models"
	"github.com/go-chi/render"
	"github.com/xuri/excelize/v2"
)

type BulkTasksOperations interface {
	Create(t []models.Task) (tasks []models.Task, err error)
}

type BulkTasksHandlers struct {
	Operations BulkTasksOperations
}

func NewBulkTasksHandlers(operations BulkTasksOperations) *BulkTasksHandlers {
	return &BulkTasksHandlers{
		Operations: operations,
	}
}

func (h *BulkTasksHandlers) CreateHandler(w http.ResponseWriter, r *http.Request) {
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

	t := []models.Task{}
	for i, row := range rows {
		if i == 0 {
			continue
		}

		task := models.Task{
			Name:   row[0],
			Status: row[1],
		}
		t = append(t, task)
	}

	tasks, err := h.Operations.Create(t)
	if err != nil {
		render.Render(w, r, errRender(err))
		return
	}

	render.Respond(w, r, newTaskListResponse(&tasks))
}
