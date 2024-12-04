package handlers

import (
	"fmt"

	"github.com/gatsu420/ngetes/models"
)

type tokenResponse struct {
	Token string `json:"access_token"`
}

func newTokenResponse(token string) *tokenResponse {
	return &tokenResponse{
		Token: token,
	}
}

type tokenClaimResponse struct {
	Claim map[string]interface{} `json:"claim"`
}

func newTokenClaimResponse(c map[string]interface{}) *tokenClaimResponse {
	return &tokenClaimResponse{
		Claim: c,
	}
}

type taskListResponse struct {
	Task *[]models.Task `json:"tasks"`
}

func newTaskListResponse(t *[]models.Task) *taskListResponse {
	return &taskListResponse{
		Task: t,
	}
}

type taskResponse struct {
	Task *models.Task `json:"task"`
}

func newTaskResponse(t *models.Task) *taskResponse {
	return &taskResponse{
		Task: t,
	}
}

type deletedTaskResponse struct {
	Status string `json:"status"`
}

func newDeletedTaskResponse(t *models.Task) *deletedTaskResponse {
	return &deletedTaskResponse{
		Status: fmt.Sprintf("deleted task ID %v", t.ID),
	}
}

type userResponse struct {
	User *models.User `json:"user"`
}

func newUserResponse(u *models.User) *userResponse {
	return &userResponse{
		User: u,
	}
}
