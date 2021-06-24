package store

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddPost(t *testing.T) {
	testSetup()
	user, err := addTestUser()
	assert.NoError(t, err)

	post, err := addTestPost(user)
	assert.NoError(t, err)
	assert.Equal(t, 1, post.ID)
}

func TestFetchUserPosts(t *testing.T) {
	testSetup()
	user, err := addTestUser()
	assert.NoError(t, err)
	post, err := addTestPost(user)
	assert.NoError(t, err)

	err = FetchUserPosts(user)
	assert.NoError(t, err)
	assert.Equal(t, post, user.Posts[0])
}

func TestFetchUserPostsEmpty(t *testing.T) {
	testSetup()
	user, err := addTestUser()
	assert.NoError(t, err)

	err = FetchUserPosts(user)
	assert.NoError(t, err)
	assert.Empty(t, user.Posts)
	assert.NotNil(t, user.Posts)
}

func TestFetchPost(t *testing.T) {
	testSetup()
	user, err := addTestUser()
	assert.NoError(t, err)
	post, err := addTestPost(user)
	assert.NoError(t, err)

	fetchedPost, err := FetchPost(post.ID)
	assert.NoError(t, err)
	assert.Equal(t, post.ID, fetchedPost.ID)
	assert.Equal(t, post.Title, fetchedPost.Title)
	assert.Equal(t, post.Content, fetchedPost.Content)
	assert.Equal(t, user.ID, fetchedPost.UserID)
}

func TestFetchNotExistingPost(t *testing.T) {
	testSetup()

	fetchedPost, err := FetchPost(1)
	assert.Error(t, err)
	assert.Nil(t, fetchedPost)
	assert.Equal(t, "Not found.", err.Error())
}

func TestUpdatePost(t *testing.T) {
	testSetup()
	user, err := addTestUser()
	assert.NoError(t, err)
	post, err := addTestPost(user)
	assert.NoError(t, err)

	post.Title = "New title"
	post.Content = "New content"
	err = UpdatePost(post)
	assert.NoError(t, err)
}

func TestDeletePost(t *testing.T) {
	testSetup()
	user, err := addTestUser()
	assert.NoError(t, err)
	post, err := addTestPost(user)
	assert.NoError(t, err)

	err = DeletePost(post)
	assert.NoError(t, err)
}
