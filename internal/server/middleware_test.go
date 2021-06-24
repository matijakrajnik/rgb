package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuthorizationHeaderInvalidFormat(t *testing.T) {
	router := testSetup()
	user := addTestUser()
	token := generateJWT(user)

	req := NewRequest(router, "GET", "/api/posts", "")
	rec := httptest.NewRecorder()
	req.Header.Add("Authorization", "Bearer"+token)
	router.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	assert.Equal(t, "Authorization header format is not valid.", jsonRes(rec.Body)["error"])
}

func TestAuthorizationHeaderMissingBearer(t *testing.T) {
	router := testSetup()
	user := addTestUser()
	token := generateJWT(user)

	req := NewRequest(router, "GET", "/api/posts", "")
	rec := httptest.NewRecorder()
	req.Header.Add("Authorization", "Bearr "+token)
	router.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	assert.Equal(t, "Authorization header is missing bearer part.", jsonRes(rec.Body)["error"])
}

func TestAuthorizationInvalidToken(t *testing.T) {
	router := testSetup()
	user := addTestUser()
	token := generateJWT(user)

	req := NewRequest(router, "GET", "/api/posts", "")
	rec := httptest.NewRecorder()
	req.Header.Add("Authorization", "Bearer invalid"+token)
	router.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	assert.Equal(t, "jwt: token format is not valid", jsonRes(rec.Body)["error"])
}
