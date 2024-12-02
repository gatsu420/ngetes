package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/jwtauth/v5"
	"github.com/go-chi/render"
)

type AuthOperations interface {
	CreateJWTAuth() (*jwtauth.JWTAuth, error)
}

type AuthHandlers struct {
	Operations     AuthOperations
	UserOperations UserOperations
}

func NewAuthHandlers(operations AuthOperations, userOperations UserOperations) *AuthHandlers {
	return &AuthHandlers{
		Operations:     operations,
		UserOperations: userOperations,
	}
}

type tokenResponse struct {
	Token string `json:"access_token"`
}

func newTokenResponse(token string) *tokenResponse {
	return &tokenResponse{
		Token: token,
	}
}

type validUserNameResponse struct {
	Response string `json:"message"`
}

func newValidUserNameResponse(userName string) *validUserNameResponse {
	return &validUserNameResponse{
		Response: fmt.Sprintf("welcome, %v", userName),
	}
}

func (hd *AuthHandlers) GetTokenHandler(w http.ResponseWriter, r *http.Request) {
	jwtAuth, err := hd.Operations.CreateJWTAuth()
	if err != nil {
		render.Render(w, r, errRender(err))
		return
	}

	_, token, err := jwtAuth.Encode(map[string]interface{}{
		"user_id": "ngetes doankk",
		"exp":     time.Now().Add(100 * time.Second).Unix(),
	})
	if err != nil {
		render.Render(w, r, errRender(err))
		return
	}

	render.Respond(w, r, newTokenResponse(token))
}

func (hd *AuthHandlers) GetValidUserNameHandler(w http.ResponseWriter, r *http.Request) {
	user := &userRequest{}
	err := render.Bind(r, user)
	if err != nil {
		render.Render(w, r, errRender(err))
		return
	}

	existence, err := hd.UserOperations.GetValidUserName(user.User, user.User.Name)
	if err != nil {
		render.Render(w, r, errRender(err))
		return
	}

	if !existence {
		render.Render(w, r, errUnauthorizedRender())
	} else {
		render.Respond(w, r, newValidUserNameResponse(user.User.Name))
	}
}
