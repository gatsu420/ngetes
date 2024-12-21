package auth_test

import (
	"reflect"
	"testing"

	"github.com/gatsu420/ngetes/auth"
	"github.com/go-chi/jwtauth/v5"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestNewAuthStore(t *testing.T) {
	tests := []struct {
		testName           string
		jwtauth            *jwtauth.JWTAuth
		cache              *redis.Client
		expectedStoreIsNil bool
		expectedStoreType  reflect.Type
	}{
		{
			testName:           "valid store",
			jwtauth:            jwtauth.New("HS256", []byte("hahaha"), nil),
			cache:              redis.NewClient(&redis.Options{}),
			expectedStoreIsNil: false,
			expectedStoreType:  reflect.TypeOf(&auth.AuthStore{}),
		},
		{
			testName:           "empty store",
			jwtauth:            nil,
			cache:              nil,
			expectedStoreIsNil: false,
			expectedStoreType:  reflect.TypeOf(&auth.AuthStore{}),
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			authStore := auth.NewAuthStore(test.jwtauth, test.cache)

			assert.Equal(t, test.expectedStoreIsNil, authStore == nil, "wrong store state")
			assert.Equal(t, test.expectedStoreType, reflect.TypeOf(authStore), "wrong store type")
		})
	}
}
