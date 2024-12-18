package auth_test

import (
	"reflect"
	"testing"

	"github.com/gatsu420/ngetes/auth"
	"github.com/gatsu420/ngetes/config"
	"github.com/go-chi/jwtauth/v5"
	"github.com/stretchr/testify/assert"
)

func TestJWTAuth(t *testing.T) {
	tests := []struct {
		testName          string
		config            *config.Config
		expectedAuthIsNil bool
		expectedAuthType  reflect.Type
		expectedErrIsNil  bool
	}{
		{
			testName:          "valid auth",
			config:            &config.Config{TokenSecretKey: "hahaha"},
			expectedAuthIsNil: false,
			expectedAuthType:  reflect.TypeOf(&jwtauth.JWTAuth{}),
			expectedErrIsNil:  true,
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			auth, err := auth.JWTAuth(test.config)

			assert.Equal(t, test.expectedAuthIsNil, auth == nil, "wrong auth state")
			assert.Equal(t, test.expectedAuthType, reflect.TypeOf(auth), "wrong auth type")
			assert.Equal(t, test.expectedErrIsNil, err == nil, "wrong error state")
		})
	}
}
