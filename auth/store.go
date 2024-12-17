package auth

import (
	"github.com/go-chi/jwtauth/v5"
	"github.com/redis/go-redis/v9"
)

type AuthStore struct {
	jwtAuth *jwtauth.JWTAuth
	cache   *redis.Client
}

func NewAuthStore(jwtauth *jwtauth.JWTAuth, cache *redis.Client) *AuthStore {
	return &AuthStore{
		jwtAuth: jwtauth,
		cache:   cache,
	}
}
