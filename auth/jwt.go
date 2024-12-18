package auth

import (
	"github.com/gatsu420/ngetes/config"
	"github.com/go-chi/jwtauth/v5"
)

func JWTAuth(config *config.Config) (*jwtauth.JWTAuth, error) {
	auth := jwtauth.New("HS256", []byte(config.TokenSecretKey), nil)

	return auth, nil
}
