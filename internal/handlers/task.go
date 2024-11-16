package handlers

import (
	"net/http"

	"github.com/gatsu420/ngetes/models"
)

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
