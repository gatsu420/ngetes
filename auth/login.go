package auth

import (
	"github.com/go-chi/jwtauth/v5"
)

func (s *authStore) CreateJWTAuth() (*jwtauth.JWTAuth, error) {
	auth, err := JWTAuth()

	return auth, err
}
