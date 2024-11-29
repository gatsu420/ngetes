package auth

import (
	"github.com/go-chi/jwtauth/v5"
)

func (s *AuthStore) CreateJWTAuth() (*jwtauth.JWTAuth, error) {
	auth, err := JWTAuth()

	return auth, err
}
