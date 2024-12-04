package handlers

import (
	"net/http"

	"github.com/gatsu420/ngetes/models"
)

type taskRequest struct {
	Task *models.Task `json:"task"`
}

func (tr *taskRequest) Bind(r *http.Request) error {
	return nil
}

type userRequest struct {
	User *models.User `json:"user"`
}

func (ur *userRequest) Bind(r *http.Request) error {
	return nil
}
