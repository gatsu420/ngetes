package handlers

import (
	"net/http"
	"time"

	"github.com/go-chi/jwtauth/v5"
	"github.com/go-chi/render"
)

type AuthOperations interface {
	CreateJWTAuth() (*jwtauth.JWTAuth, error)
	GetJWTClaim(r *http.Request) (map[string]interface{}, error)
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

func (hd *AuthHandlers) GetTokenHandler(w http.ResponseWriter, r *http.Request) {
	user := &userRequest{}
	err := render.Bind(r, user)
	if err != nil {
		render.Render(w, r, errRender(err))
		return
	}

	existence, err := hd.UserOperations.GetValidUserName(user.User.Name)
	if err != nil {
		render.Render(w, r, errRender(err))
		return
	}

	if !existence {
		render.Render(w, r, errUnauthorizedRender())
		return
	}

	jwtAuth, err := hd.Operations.CreateJWTAuth()
	if err != nil {
		render.Render(w, r, errRender(err))
		return
	}

	_, token, err := jwtAuth.Encode(map[string]interface{}{
		"exp":      time.Now().Add(100 * time.Second).Unix(),
		"userName": user.User.Name,
	})
	if err != nil {
		render.Render(w, r, errRender(err))
		return
	}

	render.Respond(w, r, newTokenResponse(token))
}

func (hd *AuthHandlers) GetTokenClaimHandler(w http.ResponseWriter, r *http.Request) {
	claim := r.Context().Value(authTokenClaimCtx).(map[string]interface{})

	render.Respond(w, r, newTokenClaimResponse(claim))
}
