package handlers

import (
	"net/http"
	"time"

	"github.com/gatsu420/ngetes/auth"
	"github.com/go-chi/render"
)

type tokenResponse struct {
	Token string `json:"access_token"`
}

func newTokenResponse(token string) *tokenResponse {
	return &tokenResponse{
		Token: token,
	}
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	auth := auth.NewAuth()

	_, token, err := auth.Encode(map[string]interface{}{
		"user_id": "ngetes",
		"exp":     time.Now().Add(1120 * time.Second).Unix(),
	})
	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	render.Respond(w, r, newTokenResponse(token))
}
