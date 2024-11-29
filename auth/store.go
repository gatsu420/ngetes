package auth

import "github.com/go-chi/jwtauth/v5"

type authStore struct {
	jwtAuth *jwtauth.JWTAuth
}

func NewAuthStore(jwtauth *jwtauth.JWTAuth) *authStore {
	return &authStore{
		jwtAuth: jwtauth,
	}
}
