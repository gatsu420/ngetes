package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type ctxKey int

const (
	authTokenClaimCtx ctxKey = iota
	taskCtx
)

func (hd *AuthHandlers) TokenClaimCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claim, err := hd.Operations.GetJWTClaim(r)
		if err != nil {
			render.Render(w, r, errRender(err))
			return
		}

		ctx := context.WithValue(r.Context(), authTokenClaimCtx, claim)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (hd *AuthHandlers) TokenBlacklistAccess(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claim := r.Context().Value(authTokenClaimCtx).(map[string]interface{})
		userName := claim["userName"].(string)

		isTokenBlacklisted, err := hd.Operations.GetTokenBlacklistFlag(userName)
		if err != nil {
			render.Render(w, r, errRender(err))
			return
		}

		if isTokenBlacklisted {
			render.Render(w, r, errUnauthorizedRender())
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (hd *AuthHandlers) AdminAccess(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claim := r.Context().Value(authTokenClaimCtx).(map[string]interface{})
		userName := claim["userName"].(string)

		roleID, err := hd.UserOperations.GetRoleByUserName(userName)
		if err != nil {
			render.Render(w, r, errRender(err))
			return
		}
		if roleID != 1 {
			render.Render(w, r, errForbiddenRender())
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (hd *TaskHandlers) TaskCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "taskID"))
		if err != nil {
			render.Render(w, r, errRender(err))
			return
		}

		task, err := hd.Operations.Get(id)
		if err != nil {
			render.Render(w, r, errRender(err))
			return
		}

		ctx := context.WithValue(r.Context(), taskCtx, task)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
