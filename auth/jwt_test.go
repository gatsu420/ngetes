package auth_test

import (
	"testing"

	"github.com/gatsu420/ngetes/auth"
	"github.com/gatsu420/ngetes/config"
	"github.com/stretchr/testify/assert"
)

func TestJWTAuth(t *testing.T) {
	mockConfig := &config.Config{
		TokenSecretKey: "hahaha",
	}

	auth, err := auth.JWTAuth(mockConfig)

	assert.NoError(t, err, "must be no error")
	assert.NotNil(t, auth, "must be instantiated as non nil")
}
