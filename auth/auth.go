package auth

import (
	"context"
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

func (s *AuthStore) CreateUserMemory(userName string) error {
	ctx := context.Background()

	_, err := s.RDB.JSONSet(ctx, userName, "$", `{"isTokenBlacklisted": false}`).Result()
	if err != nil {
		return err
	}

	return nil
}

func (s *AuthStore) UpdateTokenBlacklistFlag(userName string, isBlacklisted bool) error {
	ctx := context.Background()

	_, err := s.RDB.JSONSet(ctx, userName, ".isTokenBlacklisted", isBlacklisted).Result()
	if err != nil {
		return err
	}

	return nil
}

func (s *AuthStore) GetTokenExistence(userName string) (isExist bool, err error) {
	ctx := context.Background()

	val, err := s.RDB.JSONGet(ctx, userName, "$").Result()
	if err != nil {
		return false, err
	}

	return val != "", nil
}

func (s *AuthStore) GetTokenBlacklistFlag(userName string) (flag string, err error) {
	ctx := context.Background()

	val, err := s.RDB.JSONGet(ctx, userName, ".isTokenBlacklisted").Result()
	if err != nil {
		return "", err
	}

	return val, nil
}
