package handlers

import (
	"net/http"
	"time"

	"github.com/go-chi/jwtauth/v5"
	"github.com/go-chi/render"
)

type AuthOperations interface {
	CreateJWTAuth() (*jwtauth.JWTAuth, error)
}

type AuthHandlers struct {
	Operations AuthOperations
}

func NewAuthHandlers(operations AuthOperations) *AuthHandlers {
	return &AuthHandlers{
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

func (hd *AuthHandlers) LoginHandler(w http.ResponseWriter, r *http.Request) {
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
