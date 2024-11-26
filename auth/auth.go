package auth

import "github.com/go-chi/jwtauth/v5"

type authStore struct {
	JWTAuth *jwtauth.JWTAuth
}

func NewAuthStore(jwtauth *jwtauth.JWTAuth) *authStore {
	return &authStore{
		JWTAuth: jwtauth,
	}
}
