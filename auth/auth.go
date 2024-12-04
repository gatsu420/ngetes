package auth

import (
	"net/http"

	"github.com/go-chi/jwtauth/v5"
)

func (s *AuthStore) CreateJWTAuth() (*jwtauth.JWTAuth, error) {
	auth, err := JWTAuth()

	return auth, err
}

func (s *AuthStore) GetJWTClaim(r *http.Request) (map[string]interface{}, error) {
	_, claim, err := jwtauth.FromContext(r.Context())
	if err != nil {
		return nil, err
	}

	return claim, nil
}
