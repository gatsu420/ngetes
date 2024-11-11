package handlers

import (
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
