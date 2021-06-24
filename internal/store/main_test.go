package store

import "github.com/gin-gonic/gin"

func testSetup() {
	gin.SetMode(gin.TestMode)
	ResetTestDatabase()
}

func addTestUser() (*User, error) {
	user := &User{
		Username: "batman",
		Password: "secret123",
	}
	err := AddUser(user)
	return user, err
}

func addTestPost(user *User) (*Post, error) {
	post := &Post{
		Title:   "Gotham cronicles",
		Content: "Joker is planning big hit tonight.",
	}
	err := AddPost(user, post)
	return post, err
}
