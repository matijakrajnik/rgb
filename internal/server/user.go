package server

import (
	"net/http"
	"rgb/internal/store"

	"github.com/gin-gonic/gin"
)

func signUp(ctx *gin.Context) {
	user := ctx.MustGet(gin.BindKey).(*store.User)
	if err := store.AddUser(user); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "Signed up successfully.",
		"jwt": generateJWT(user),
	})
}

func signIn(ctx *gin.Context) {
	user := ctx.MustGet(gin.BindKey).(*store.User)
	user, err := store.Authenticate(user.Username, user.Password)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Sign in failed."})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg": "Signed in successfully.",
		"jwt": generateJWT(user),
	})
}
