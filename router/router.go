package router

import (
	"twittir-go/controllers"
	"twittir-go/middleware"

	"github.com/gin-gonic/gin"
)

func StartServer() *gin.Engine {
	router := gin.Default()

	apiRouter := router.Group("/api")
	{
		// Auth
		apiRouter.POST("/register", controllers.UserRegister)
		apiRouter.POST("/login", controllers.UserLogin)
		apiRouter.GET("/profile", middleware.Authentication(), controllers.GetDetailUser)

		// Post
		apiRouter.POST("/posts", middleware.Authentication(), controllers.CreatePost)
		apiRouter.GET("/posts", middleware.Authentication(), controllers.GetPosts)

		// Comment post
		apiRouter.POST("/posts/:id/comment", middleware.Authentication(), controllers.CreateComment)

		// Like Post
		apiRouter.POST("/posts/:id/like", middleware.Authentication(), controllers.CreateLikePost)
		apiRouter.POST("/posts/:id/comment/like", middleware.Authentication(), controllers.CreateLikeComment)
	}
	return router
}
