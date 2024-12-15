package auth

import (
	"github.com/gatsu420/ngetes/logger"
	"github.com/go-chi/jwtauth/v5"
	"github.com/spf13/viper"
	"go.uber.org/zap"
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
		logger.Logger.Error("failed to read config file", zap.Error(err))
	}

	tokenConfig.secretKey = viper.GetString("TOKEN_SECRET_KEY")
}

func JWTAuth() (*jwtauth.JWTAuth, error) {
	auth := jwtauth.New("HS256", []byte(tokenConfig.secretKey), nil)

	return auth, nil
}
