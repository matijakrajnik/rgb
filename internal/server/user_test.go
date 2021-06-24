package server

import (
	"net/http"
	"rgb/internal/store"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSignUp(t *testing.T) {
	router := testSetup()

	body := userJSON(store.User{
		Username: "batman",
		Password: "secret123",
	})
	rec := performRequest(router, "POST", "/api/signup", body)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "Signed up successfully.", jsonRes(rec.Body)["msg"])
	assert.NotEmpty(t, jsonRes(rec.Body)["jwt"])
}

func TestSignUpEmptyUsername(t *testing.T) {
	router := testSetup()

	body := userJSON(store.User{
		Username: "",
		Password: "secret123",
	})
	rec := performRequest(router, "POST", "/api/signup", body)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, "Username is required.", jsonFieldError(jsonRes(rec.Body), "Username"))
	assert.Empty(t, jsonRes(rec.Body)["jwt"])
}

func TestSignUpShortUsername(t *testing.T) {
	router := testSetup()

	body := userJSON(store.User{
		Username: "batm",
		Password: "secret123",
	})
	rec := performRequest(router, "POST", "/api/signup", body)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, "Username must be longer than or equal 5 characters.", jsonFieldError(jsonRes(rec.Body), "Username"))
	assert.Empty(t, jsonRes(rec.Body)["jwt"])
}

func TestSignUpLongUsername(t *testing.T) {
	router := testSetup()

	body := userJSON(store.User{
		Username: strings.Repeat("b", 31),
		Password: "secret123",
	})
	rec := performRequest(router, "POST", "/api/signup", body)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, "Username cannot be longer than 30 characters.", jsonFieldError(jsonRes(rec.Body), "Username"))
	assert.Empty(t, jsonRes(rec.Body)["jwt"])
}

func TestSignUpExistingUsername(t *testing.T) {
	router := testSetup()
	user := addTestUser()

	body := userJSON(store.User{
		Username: user.Username,
		Password: user.Password,
	})
	rec := performRequest(router, "POST", "/api/signup", body)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, "Username already exists.", jsonRes(rec.Body)["error"])
	assert.Empty(t, jsonRes(rec.Body)["jwt"])
}

func TestSignUpEmptyPassword(t *testing.T) {
	router := testSetup()

	body := userJSON(store.User{
		Username: "batman",
		Password: "",
	})
	rec := performRequest(router, "POST", "/api/signup", body)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, "Password is required.", jsonFieldError(jsonRes(rec.Body), "Password"))
	assert.Empty(t, jsonRes(rec.Body)["jwt"])
}

func TestSignUpShortPassword(t *testing.T) {
	router := testSetup()

	body := userJSON(store.User{
		Username: "batman",
		Password: "secret",
	})
	rec := performRequest(router, "POST", "/api/signup", body)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, "Password must be longer than or equal 7 characters.", jsonFieldError(jsonRes(rec.Body), "Password"))
	assert.Empty(t, jsonRes(rec.Body)["jwt"])
}

func TestSignUpLongPassword(t *testing.T) {
	router := testSetup()

	body := userJSON(store.User{
		Username: "batman",
		Password: strings.Repeat("s", 33),
	})
	rec := performRequest(router, "POST", "/api/signup", body)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, "Password cannot be longer than 32 characters.", jsonFieldError(jsonRes(rec.Body), "Password"))
	assert.Empty(t, jsonRes(rec.Body)["jwt"])
}

func TestSignIn(t *testing.T) {
	router := testSetup()
	user := addTestUser()

	body := userJSON(store.User{
		Username: user.Username,
		Password: user.Password,
	})
	rec := performRequest(router, "POST", "/api/signin", body)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "Signed in successfully.", jsonRes(rec.Body)["msg"])
	assert.NotEmpty(t, jsonRes(rec.Body)["jwt"])
}

func TestSignInInvalidUsername(t *testing.T) {
	router := testSetup()
	user := addTestUser()

	body := userJSON(store.User{
		Username: "invalid",
		Password: user.Password,
	})
	rec := performRequest(router, "POST", "/api/signin", body)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	assert.Equal(t, "Sign in failed.", jsonRes(rec.Body)["error"])
	assert.Empty(t, jsonRes(rec.Body)["jwt"])
}

func TestSignInInvalidPassword(t *testing.T) {
	router := testSetup()
	user := addTestUser()

	body := userJSON(store.User{
		Username: user.Username,
		Password: "invalid",
	})
	rec := performRequest(router, "POST", "/api/signin", body)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	assert.Equal(t, "Sign in failed.", jsonRes(rec.Body)["error"])
	assert.Empty(t, jsonRes(rec.Body)["jwt"])
}
