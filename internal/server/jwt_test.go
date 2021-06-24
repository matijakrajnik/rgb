package server

import (
	"rgb/internal/conf"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJwtSetup(t *testing.T) {
	_ = testSetup()
	assert.NotPanics(t, func() { jwtSetup(conf.NewConfig("dev")) })
	assert.NotNil(t, jwtSigner)
	assert.NotNil(t, jwtVerifier)
}

func TestGenerateJWT(t *testing.T) {
	_ = testSetup()
	user := addTestUser()

	token := generateJWT(user)
	assert.NotEmpty(t, token)
}

func TestVerifyJWT(t *testing.T) {
	_ = testSetup()
	user := addTestUser()
	token := generateJWT(user)
	assert.NotEmpty(t, token)

	userID, err := verifyJWT(token)
	assert.NoError(t, err)
	assert.Equal(t, user.ID, userID)
}

func TestVerifyInvalidJWT(t *testing.T) {
	_ = testSetup()
	user := addTestUser()
	token := generateJWT(user)
	assert.NotEmpty(t, token)

	userID, err := verifyJWT(token + "invalid")
	assert.Error(t, err)
	assert.Equal(t, 0, userID)
}
