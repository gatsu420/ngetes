package auth

import (
	"github.com/go-chi/jwtauth/v5"
	"github.com/redis/go-redis/v9"
)

type AuthStore struct {
	jwtAuth *jwtauth.JWTAuth
	RDB     *redis.Client
}

func NewAuthStore(jwtauth *jwtauth.JWTAuth, rdb *redis.Client) *AuthStore {
	return &AuthStore{
		jwtAuth: jwtauth,
		RDB:     rdb,
	}
}
