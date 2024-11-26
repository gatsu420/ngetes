package auth

import (
	"fmt"

	"github.com/go-chi/jwtauth/v5"
	"github.com/spf13/viper"
)

func JWTAuth() (*jwtauth.JWTAuth, error) {
	viper.SetConfigFile("./.env")
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	secretKey := viper.GetString("token_secret_key")
	auth := jwtauth.New("HS256", []byte(secretKey), nil)

	return auth, nil
}
