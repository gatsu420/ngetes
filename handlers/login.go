package handlers

import (
	"net/http"
	"time"

	"github.com/go-chi/jwtauth/v5"
	"github.com/go-chi/render"
)

type LoginOperations interface {
	CreateJWTAuth() (*jwtauth.JWTAuth, error)
}

type LoginHandlers struct {
	Operations LoginOperations
}

func NewLoginHandlers(operations LoginOperations) *LoginHandlers {
	return &LoginHandlers{
		Operations: operations,
	}
}

type authResponse struct {
	Token string `json:"access_token"`
}

func newAuthResponse(token string) *authResponse {
	return &authResponse{
		Token: token,
	}
}

func (hd *LoginHandlers) LoginHandler(w http.ResponseWriter, r *http.Request) {
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

	render.Respond(w, r, newAuthResponse(token))
}
