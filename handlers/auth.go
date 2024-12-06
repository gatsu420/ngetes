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

	CreateUserMemory(userName string) error
	UpdateTokenBlacklistFlag(userName string, isBlacklisted bool) error
	GetTokenExistence(userName string) (isExist bool, err error)
	GetTokenBlacklistFlag(userName string) (flag string, err error)
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

	existence, err := hd.UserOperations.GetUserNameExistence(user.User.Name)
	if err != nil {
		render.Render(w, r, errRender(err))
		return
	}

	if !existence {
		render.Render(w, r, errUnauthorizedRender())
		return
	}

	// This token validation thing need to be put on middleware, instead of single endpoint
	// to make sure the effect propagates to other endpoints.
	isTokenExist, err := hd.Operations.GetTokenExistence(user.User.Name)
	if err != nil {
		render.Render(w, r, errRender(err))
		return
	}

	if !isTokenExist {
		err = hd.Operations.CreateUserMemory(user.User.Name)
		if err != nil {
			render.Render(w, r, errRender(err))
			return
		}
	} else {
		tokenBlacklistFlag, err := hd.Operations.GetTokenBlacklistFlag(user.User.Name)
		if err != nil {
			render.Render(w, r, errRender(err))
			return
		}

		if tokenBlacklistFlag == "true" {
			render.Render(w, r, errForbiddenRender())
			return
		}
	}

	err = hd.Operations.UpdateTokenBlacklistFlag(user.User.Name, false)
	if err != nil {
		render.Render(w, r, errRender(err))
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

func (hd *AuthHandlers) CreateTokenBlacklistHandler(w http.ResponseWriter, r *http.Request) {
	claim := r.Context().Value(authTokenClaimCtx).(map[string]interface{})
	userName := claim["userName"].(string)

	err := hd.Operations.UpdateTokenBlacklistFlag(userName, true)
	if err != nil {
		render.Render(w, r, errRender(err))
		return
	}

	render.Respond(w, r, "token blacklisted")
}
