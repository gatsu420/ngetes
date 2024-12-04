package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/jwtauth/v5"
	"github.com/go-chi/render"
)

type ctxAuthKey int

const ctxAuthTokenClaim ctxAuthKey = iota

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

func (hd *AuthHandlers) TokenClaimCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claim, err := hd.Operations.GetJWTClaim(r)
		if err != nil {
			render.Render(w, r, errRender(err))
			return
		}

		ctx := context.WithValue(r.Context(), ctxAuthTokenClaim, claim)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
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
	claim, err := hd.Operations.GetJWTClaim(r)
	if err != nil {
		render.Render(w, r, errRender(err))
		return
	}

	render.Respond(w, r, newTokenClaimResponse(claim))
}
