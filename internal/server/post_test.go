package server

import (
	"fmt"
	"net/http"
	"rgb/internal/store"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreatePost(t *testing.T) {
	router := testSetup()
	user := addTestUser()
	token := generateJWT(user)

	post := store.Post{
		Title:   "Gotham cronicles",
		Content: "Joker is planning big hit tonight.",
	}
	body := postJSON(post)
	rec := PerformAuthorizedRequest(router, token, "POST", "/api/posts", body)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "Post created successfully.", jsonRes(rec.Body)["msg"])
	assert.Equal(t, float64(1), jsonFieldData(jsonRes(rec.Body), "ID"))
	assert.Equal(t, post.Title, jsonFieldData(jsonRes(rec.Body), "Title"))
	assert.Equal(t, post.Content, jsonFieldData(jsonRes(rec.Body), "Content"))
	assert.NotEmpty(t, post.Content, jsonFieldData(jsonRes(rec.Body), "CreatedAt"))
	assert.NotEmpty(t, post.Content, jsonFieldData(jsonRes(rec.Body), "ModifiedAt"))
}

func TestCreatePostUnathorized(t *testing.T) {
	router := testSetup()

	post := store.Post{
		Title:   "Gotham cronicles",
		Content: "Joker is planning big hit tonight.",
	}
	body := postJSON(post)
	rec := performRequest(router, "POST", "/api/posts", body)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	assert.Equal(t, "Authorization header missing.", jsonRes(rec.Body)["error"])
}

func TestCreatePostEmptyTitle(t *testing.T) {
	router := testSetup()
	user := addTestUser()
	token := generateJWT(user)

	post := store.Post{
		Title:   "",
		Content: "Joker is planning big hit tonight.",
	}
	body := postJSON(post)
	rec := PerformAuthorizedRequest(router, token, "POST", "/api/posts", body)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, "Title is required.", jsonFieldError(jsonRes(rec.Body), "Title"))
}

func TestCreatePostShortTitle(t *testing.T) {
	router := testSetup()
	user := addTestUser()
	token := generateJWT(user)

	post := store.Post{
		Title:   "Go",
		Content: "Joker is planning big hit tonight.",
	}
	body := postJSON(post)
	rec := PerformAuthorizedRequest(router, token, "POST", "/api/posts", body)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, "Title must be longer than or equal 3 characters.", jsonFieldError(jsonRes(rec.Body), "Title"))
}

func TestCreatePostLongTitle(t *testing.T) {
	router := testSetup()
	user := addTestUser()
	token := generateJWT(user)

	post := store.Post{
		Title:   strings.Repeat("G", 51),
		Content: "Joker is planning big hit tonight.",
	}
	body := postJSON(post)
	rec := PerformAuthorizedRequest(router, token, "POST", "/api/posts", body)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, "Title cannot be longer than 50 characters.", jsonFieldError(jsonRes(rec.Body), "Title"))
}

func TestCreatePostEmptyContent(t *testing.T) {
	router := testSetup()
	user := addTestUser()
	token := generateJWT(user)

	post := store.Post{
		Title:   "Gotham cronicles",
		Content: "",
	}
	body := postJSON(post)
	rec := PerformAuthorizedRequest(router, token, "POST", "/api/posts", body)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, "Content is required.", jsonFieldError(jsonRes(rec.Body), "Content"))
}

func TestCreatePostShortContent(t *testing.T) {
	router := testSetup()
	user := addTestUser()
	token := generateJWT(user)

	post := store.Post{
		Title:   "Gotham cronicles",
		Content: "Joke",
	}
	body := postJSON(post)
	rec := PerformAuthorizedRequest(router, token, "POST", "/api/posts", body)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, "Content must be longer than or equal 5 characters.", jsonFieldError(jsonRes(rec.Body), "Content"))
}

func TestCreatePostLongContent(t *testing.T) {
	router := testSetup()
	user := addTestUser()
	token := generateJWT(user)

	post := store.Post{
		Title:   "Gotham cronicles",
		Content: strings.Repeat("J", 5001),
	}
	body := postJSON(post)
	rec := PerformAuthorizedRequest(router, token, "POST", "/api/posts", body)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, "Content cannot be longer than 5000 characters.", jsonFieldError(jsonRes(rec.Body), "Content"))
}

func TestIndexPosts(t *testing.T) {
	router := testSetup()
	user := addTestUser()
	token := generateJWT(user)
	post := addTestPost(user)

	rec := PerformAuthorizedRequest(router, token, "GET", "/api/posts", "")
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "Posts fetched successfully.", jsonRes(rec.Body)["msg"])
	assert.Equal(t, float64(post.ID), jsonDataSlice(rec.Body)[0]["ID"])
	assert.Equal(t, post.Title, jsonDataSlice(rec.Body)[0]["Title"])
	assert.Equal(t, post.Content, jsonDataSlice(rec.Body)[0]["Content"])
	assert.NotEmpty(t, post.Content, jsonDataSlice(rec.Body)[0]["CreatedAt"])
	assert.NotEmpty(t, post.Content, jsonDataSlice(rec.Body)[0]["ModifiedAt"])
}

func TestIndexPostsUnathorized(t *testing.T) {
	router := testSetup()
	user := addTestUser()
	_ = addTestPost(user)

	rec := performRequest(router, "GET", "/api/posts", "")
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	assert.Equal(t, "Authorization header missing.", jsonRes(rec.Body)["error"])
}

func TestIndexPostOnlyOwned(t *testing.T) {
	router := testSetup()
	user1 := addTestUser()
	user2 := addTestUser2()
	token1 := generateJWT(user1)
	token2 := generateJWT(user2)
	post1 := addTestPost(user1)
	post2 := addTestPost2(user2)

	rec := PerformAuthorizedRequest(router, token1, "GET", "/api/posts", "")
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "Posts fetched successfully.", jsonRes(rec.Body)["msg"])
	assert.Len(t, jsonDataSlice(rec.Body), 1)
	assert.Equal(t, float64(post1.ID), jsonDataSlice(rec.Body)[0]["ID"])
	assert.Equal(t, post1.Title, jsonDataSlice(rec.Body)[0]["Title"])
	assert.Equal(t, post1.Content, jsonDataSlice(rec.Body)[0]["Content"])
	assert.NotEmpty(t, post1.Content, jsonDataSlice(rec.Body)[0]["CreatedAt"])
	assert.NotEmpty(t, post1.Content, jsonDataSlice(rec.Body)[0]["ModifiedAt"])

	rec = PerformAuthorizedRequest(router, token2, "GET", "/api/posts", "")
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "Posts fetched successfully.", jsonRes(rec.Body)["msg"])
	assert.Len(t, jsonDataSlice(rec.Body), 1)
	assert.Equal(t, float64(post2.ID), jsonDataSlice(rec.Body)[0]["ID"])
	assert.Equal(t, post2.Title, jsonDataSlice(rec.Body)[0]["Title"])
	assert.Equal(t, post2.Content, jsonDataSlice(rec.Body)[0]["Content"])
	assert.NotEmpty(t, post2.Content, jsonDataSlice(rec.Body)[0]["CreatedAt"])
	assert.NotEmpty(t, post2.Content, jsonDataSlice(rec.Body)[0]["ModifiedAt"])
}

func TestIndexPostsEmpty(t *testing.T) {
	router := testSetup()
	user := addTestUser()
	token := generateJWT(user)

	rec := PerformAuthorizedRequest(router, token, "GET", "/api/posts", "")
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Empty(t, jsonRes(rec.Body)["data"])
	assert.NotNil(t, jsonRes(rec.Body)["data"])
	assert.Equal(t, "Posts fetched successfully.", jsonRes(rec.Body)["msg"])
}

func TestUpdatePost(t *testing.T) {
	router := testSetup()
	user := addTestUser()
	token := generateJWT(user)
	post := addTestPost(user)

	updated := store.Post{
		ID:      post.ID,
		Title:   "Gotham at night",
		Content: "Gotham never sleeps.",
	}
	rec := PerformAuthorizedRequest(router, token, "PUT", "/api/posts", postJSON(updated))
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "Post updated successfully.", jsonRes(rec.Body)["msg"])
}

func TestUpdatePostUnauthorized(t *testing.T) {
	router := testSetup()
	user := addTestUser()
	post := addTestPost(user)

	updated := store.Post{
		ID:      post.ID,
		Title:   "Gotham at night",
		Content: "Gotham never sleeps.",
	}
	rec := performRequest(router, "PUT", "/api/posts", postJSON(updated))
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	assert.Equal(t, "Authorization header missing.", jsonRes(rec.Body)["error"])
}

func TestUpdatePostEmptyTitle(t *testing.T) {
	router := testSetup()
	user := addTestUser()
	token := generateJWT(user)
	post := addTestPost(user)

	updated := store.Post{
		ID:      post.ID,
		Title:   "",
		Content: "Gotham never sleeps.",
	}
	rec := PerformAuthorizedRequest(router, token, "PUT", "/api/posts", postJSON(updated))
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, "Title is required.", jsonFieldError(jsonRes(rec.Body), "Title"))
}

func TestUpdatePostShortTitle(t *testing.T) {
	router := testSetup()
	user := addTestUser()
	token := generateJWT(user)
	post := addTestPost(user)

	updated := store.Post{
		ID:      post.ID,
		Title:   "Go",
		Content: "Gotham never sleeps.",
	}
	rec := PerformAuthorizedRequest(router, token, "PUT", "/api/posts", postJSON(updated))
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, "Title must be longer than or equal 3 characters.", jsonFieldError(jsonRes(rec.Body), "Title"))
}

func TestUpdatePostLongTitle(t *testing.T) {
	router := testSetup()
	user := addTestUser()
	token := generateJWT(user)
	post := addTestPost(user)

	updated := store.Post{
		ID:      post.ID,
		Title:   strings.Repeat("G", 51),
		Content: "Gotham never sleeps.",
	}
	rec := PerformAuthorizedRequest(router, token, "PUT", "/api/posts", postJSON(updated))
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, "Title cannot be longer than 50 characters.", jsonFieldError(jsonRes(rec.Body), "Title"))
}

func TestUpdatePostEmptyContent(t *testing.T) {
	router := testSetup()
	user := addTestUser()
	token := generateJWT(user)
	post := addTestPost(user)

	updated := store.Post{
		ID:      post.ID,
		Title:   "Gotham at night",
		Content: "",
	}
	rec := PerformAuthorizedRequest(router, token, "PUT", "/api/posts", postJSON(updated))
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, "Content is required.", jsonFieldError(jsonRes(rec.Body), "Content"))
}

func TestUpdatePostShortContent(t *testing.T) {
	router := testSetup()
	user := addTestUser()
	token := generateJWT(user)
	post := addTestPost(user)

	updated := store.Post{
		ID:      post.ID,
		Title:   "Gotham at night",
		Content: "Goth",
	}
	rec := PerformAuthorizedRequest(router, token, "PUT", "/api/posts", postJSON(updated))
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, "Content must be longer than or equal 5 characters.", jsonFieldError(jsonRes(rec.Body), "Content"))
}

func TestUpdatePostLongContent(t *testing.T) {
	router := testSetup()
	user := addTestUser()
	token := generateJWT(user)
	post := addTestPost(user)

	updated := store.Post{
		ID:      post.ID,
		Title:   "Gotham at night",
		Content: strings.Repeat("G", 5001),
	}
	rec := PerformAuthorizedRequest(router, token, "PUT", "/api/posts", postJSON(updated))
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, "Content cannot be longer than 5000 characters.", jsonFieldError(jsonRes(rec.Body), "Content"))
}

func TestUpdateNotOwnedPost(t *testing.T) {
	router := testSetup()
	user1 := addTestUser()
	user2 := addTestUser2()
	token1 := generateJWT(user1)
	token2 := generateJWT(user2)
	post1 := addTestPost(user1)
	post2 := addTestPost2(user2)
	updated1 := store.Post{
		ID:      post1.ID,
		Title:   "Gotham at night",
		Content: "Gotham never sleeps.",
	}
	updated2 := store.Post{
		ID:      post2.ID,
		Title:   "Lex",
		Content: "Lex has build new underground lab.",
	}

	rec := PerformAuthorizedRequest(router, token1, "PUT", "/api/posts", postJSON(updated2))
	assert.Equal(t, http.StatusForbidden, rec.Code)
	assert.Equal(t, "Not authorized.", jsonRes(rec.Body)["error"])

	rec = PerformAuthorizedRequest(router, token2, "PUT", "/api/posts", postJSON(updated1))
	assert.Equal(t, http.StatusForbidden, rec.Code)
	assert.Equal(t, "Not authorized.", jsonRes(rec.Body)["error"])
}

func TestUpdateNotExistingPost(t *testing.T) {
	router := testSetup()
	user := addTestUser()
	token := generateJWT(user)
	_ = addTestPost(user)

	updated := store.Post{
		ID:      123,
		Title:   "Gotham at night",
		Content: "Gotham never sleeps.",
	}
	rec := PerformAuthorizedRequest(router, token, "PUT", "/api/posts", postJSON(updated))
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, "Not found.", jsonRes(rec.Body)["error"])
}

func TestDeletePost(t *testing.T) {
	router := testSetup()
	user := addTestUser()
	token := generateJWT(user)
	post := addTestPost(user)

	rec := PerformAuthorizedRequest(router, token, "DELETE", fmt.Sprintf("/api/posts/%d", post.ID), "")
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "Post deleted successfully.", jsonRes(rec.Body)["msg"])
}

func TestDeletePostUnauthorized(t *testing.T) {
	router := testSetup()
	user := addTestUser()
	post := addTestPost(user)

	rec := performRequest(router, "DELETE", fmt.Sprintf("/api/posts/%d", post.ID), "")
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	assert.Equal(t, "Authorization header missing.", jsonRes(rec.Body)["error"])
}

func TestDeleteNotExistingPost(t *testing.T) {
	router := testSetup()
	user := addTestUser()
	token := generateJWT(user)

	rec := PerformAuthorizedRequest(router, token, "DELETE", "/api/posts/1", "")
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, "Not found.", jsonRes(rec.Body)["error"])
}

func TestDeletePostInvalidID(t *testing.T) {
	router := testSetup()
	user := addTestUser()
	token := generateJWT(user)

	rec := PerformAuthorizedRequest(router, token, "DELETE", "/api/posts/invalid", "")
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, "Not valid ID.", jsonRes(rec.Body)["error"])
}

func TestDeleteNotOwnedPost(t *testing.T) {
	router := testSetup()
	user1 := addTestUser()
	user2 := addTestUser2()
	token1 := generateJWT(user1)
	token2 := generateJWT(user2)
	post1 := addTestPost(user1)
	post2 := addTestPost2(user2)

	rec := PerformAuthorizedRequest(router, token1, "DELETE", fmt.Sprintf("/api/posts/%d", post2.ID), "")
	assert.Equal(t, http.StatusForbidden, rec.Code)
	assert.Equal(t, "Not authorized.", jsonRes(rec.Body)["error"])

	rec = PerformAuthorizedRequest(router, token2, "DELETE", fmt.Sprintf("/api/posts/%d", post1.ID), "")
	assert.Equal(t, http.StatusForbidden, rec.Code)
	assert.Equal(t, "Not authorized.", jsonRes(rec.Body)["error"])
}
