package server

import (
	"net/http"
	"rgb/internal/conf"
	"rgb/internal/store"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

func setRouter(cfg conf.Config) *gin.Engine {
	// Creates default gin router with Logger and Recovery middleware already attached
	router := gin.Default()

	// Enables automatic redirection if the current route can't be matched but a
	// handler for the path with (without) the trailing slash exists.
	router.RedirectTrailingSlash = true

	// Serve static files to frontend if server is started in production environment
	if cfg.Env == "prod" {
		router.Use(static.Serve("/", static.LocalFile("./assets/build", true)))
	}

	// Create API route group
	api := router.Group("/api")
	api.Use(customErrors)
	{
		api.POST("/signup", gin.Bind(store.User{}), signUp)
		api.POST("/signin", gin.Bind(store.User{}), signIn)
	}

	authorized := api.Group("/")
	authorized.Use(authorization)
	{
		authorized.GET("/posts", indexPosts)
		authorized.POST("/posts", gin.Bind(store.Post{}), createPost)
		authorized.PUT("/posts", gin.Bind(store.Post{}), updatePost)
		authorized.DELETE("/posts/:id", deletePost)
	}

	router.NoRoute(func(ctx *gin.Context) { ctx.JSON(http.StatusNotFound, gin.H{}) })

	return router
}
