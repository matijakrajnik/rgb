package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"rgb/internal/conf"
	"rgb/internal/store"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func testSetup() *gin.Engine {
	gin.SetMode(gin.TestMode)
	store.ResetTestDatabase()
	cfg := conf.NewConfig("dev")
	jwtSetup(cfg)
	return setRouter(cfg)
}

func addTestUser() *store.User {
	user := &store.User{
		Username: "batman",
		Password: "secret123",
	}
	err := store.AddUser(user)
	if err != nil {
		log.Panic().Err(err).Msg("Error adding test user.")
	}
	return user
}

func addTestUser2() *store.User {
	user := &store.User{
		Username: "superman",
		Password: "secret123",
	}
	err := store.AddUser(user)
	if err != nil {
		log.Panic().Err(err).Msg("Error adding test user.")
	}
	return user
}

func addTestPost(user *store.User) *store.Post {
	post := &store.Post{
		Title:   "Gotham cronicles",
		Content: "Joker is planning a big hit tonight.",
	}
	err := store.AddPost(user, post)
	if err != nil {
		log.Panic().Err(err).Msg("Error adding test post.")
	}
	return post
}

func addTestPost2(user *store.User) *store.Post {
	post := &store.Post{
		Title:   "Justice league meeting",
		Content: "Darkseid is plotting again.",
	}
	err := store.AddPost(user, post)
	if err != nil {
		log.Panic().Err(err).Msg("Error adding test post.")
	}
	return post
}

func userJSON(user store.User) string {
	body, err := json.Marshal(map[string]interface{}{
		"Username": user.Username,
		"Password": user.Password,
	})
	if err != nil {
		log.Panic().Err(err).Msg("Error marshalling JSON body.")
	}
	return string(body)
}

func postJSON(post store.Post) string {
	body, err := json.Marshal(map[string]interface{}{
		"ID":      post.ID,
		"Title":   post.Title,
		"Content": post.Content,
	})
	if err != nil {
		log.Panic().Err(err).Msg("Error marshalling JSON body.")
	}
	return string(body)
}

func jsonRes(body *bytes.Buffer) map[string]interface{} {
	jsonRes := &map[string]interface{}{}
	err := json.Unmarshal(body.Bytes(), jsonRes)
	if err != nil {
		log.Panic().Err(err).Msg("Error unmarshalling JSON body.")
	}
	return *jsonRes
}

func jsonDataSlice(body *bytes.Buffer) []map[string]interface{} {
	jsonRes := jsonRes(body)
	_jsonDataSlice, ok := jsonRes["data"].([]interface{})
	if !ok {
		log.Panic().Interface("jsonRes", jsonRes).Msg("JSON response data is not a slice.")
	}
	jsonSliceMaps := make([]map[string]interface{}, 0)
	for _, _jsonSliceMap := range _jsonDataSlice {
		jsonSliceMap, ok := _jsonSliceMap.(map[string]interface{})
		if !ok {
			log.Panic().Interface("_jsonSliceMap", _jsonSliceMap).Msg("JSON object in slice is not a map.")
		}
		jsonSliceMaps = append(jsonSliceMaps, jsonSliceMap)
	}
	return jsonSliceMaps
}

func jsonFieldError(jsonRes map[string]interface{}, field string) interface{} {
	jsonError, ok := jsonRes["error"].(map[string]interface{})
	if !ok {
		log.Panic().Interface("jsonRes", jsonRes).Msg("JSON response error is not a map.")
	}
	return jsonError[field]
}

func jsonFieldData(jsonRes map[string]interface{}, field string) interface{} {
	jsonData, ok := jsonRes["data"].(map[string]interface{})
	if !ok {
		log.Panic().Interface("jsonRes", jsonRes).Msg("JSON response data is not a map.")
	}
	return jsonData[field]
}

func NewRequest(router *gin.Engine, method, path, body string) *http.Request {
	req, err := http.NewRequest(method, path, strings.NewReader(body))
	if err != nil {
		log.Panic().Err(err).Msg("Error creating new request")
	}
	req.Header.Add("Content-Type", "application/json")
	return req
}

func performRequest(router *gin.Engine, method, path, body string) *httptest.ResponseRecorder {
	req := NewRequest(router, method, path, body)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	return rec
}

func PerformAuthorizedRequest(router *gin.Engine, token, method, path, body string) *httptest.ResponseRecorder {
	req := NewRequest(router, method, path, body)
	rec := httptest.NewRecorder()
	req.Header.Add("Authorization", "Bearer "+token)
	router.ServeHTTP(rec, req)
	return rec
}
