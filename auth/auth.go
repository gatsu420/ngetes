package auth

import (
	"github.com/go-chi/jwtauth/v5"
	"github.com/spf13/viper"
)

func NewAuth() *jwtauth.JWTAuth {
	auth := jwtauth.New("HS256", []byte(viper.GetString("token_secret_key")), nil)

	return auth
}
