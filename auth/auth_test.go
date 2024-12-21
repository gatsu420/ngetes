package auth_test

import (
	_ "embed"
	"reflect"
	"testing"

	"github.com/gatsu420/ngetes/auth"
	"github.com/go-chi/jwtauth/v5"
	"github.com/stretchr/testify/assert"
)

func TestCreateJWTAuth(t *testing.T) {
	tests := []struct {
		testName          string
		expectedAuthIsNil bool
		expectedAuthType  reflect.Type
		expectedErrIsNil  bool
	}{
		{
			testName:          "create valid auth",
			expectedAuthIsNil: false,
			expectedAuthType:  reflect.TypeOf(&jwtauth.JWTAuth{}),
			expectedErrIsNil:  true,
		},
	}

	for _, test := range tests {
		authStore := &auth.AuthStore{}
		jwtAuth, err := authStore.CreateJWTAuth()

		assert.Equal(t, test.expectedAuthIsNil, jwtAuth == nil, "wrong auth state")
		assert.Equal(t, test.expectedAuthType, reflect.TypeOf(jwtAuth), "wrong auth type")
		assert.Equal(t, test.expectedErrIsNil, err == nil, "wrong err state")
	}
}
