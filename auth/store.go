package auth

import "github.com/go-chi/jwtauth/v5"

type AuthStore struct {
	jwtAuth *jwtauth.JWTAuth
}

func NewAuthStore(jwtauth *jwtauth.JWTAuth) *AuthStore {
	return &AuthStore{
		jwtAuth: jwtauth,
	}
}
