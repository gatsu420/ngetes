package auth

import (
	"fmt"

	"github.com/go-chi/jwtauth/v5"
	"github.com/spf13/viper"
)

var (
	tokenConfig struct {
		secretKey string
	}
)

func init() {
	viper.SetConfigFile("./.env")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error while trying to read config file: %v", err))
	}

	tokenConfig.secretKey = viper.GetString("TOKEN_SECRET_KEY")
}

func JWTAuth() (*jwtauth.JWTAuth, error) {
	auth := jwtauth.New("HS256", []byte(tokenConfig.secretKey), nil)

	return auth, nil
}
