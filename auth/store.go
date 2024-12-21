package auth

import (
	"github.com/go-chi/jwtauth/v5"
	"github.com/redis/go-redis/v9"
)

type AuthStore struct {
	JWTAuth *jwtauth.JWTAuth
	Cache   *redis.Client
}

func NewAuthStore(jwtauth *jwtauth.JWTAuth, cache *redis.Client) *AuthStore {
	return &AuthStore{
		JWTAuth: jwtauth,
		Cache:   cache,
	}
}
