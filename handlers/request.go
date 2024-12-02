package handlers

import (
	"net/http"

	"github.com/gatsu420/ngetes/models"
)

type userRequest struct {
	User *models.User `json:"user"`
}

func (ur *userRequest) Bind(r *http.Request) error {
	return nil
}
