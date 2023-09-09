package router

import (
	"twittir-go/controllers"
	"twittir-go/middleware"

	"github.com/gin-gonic/gin"
)

func StartServer() *gin.Engine {
	router := gin.Default()

	userRouter := router.Group("/user")
	{
		userRouter.POST("/register", controllers.UserRegister)
		userRouter.POST("/login", controllers.UserLogin)
		userRouter.GET("/detail", middleware.Authentication(), controllers.GetDetailUser)
	}
	postRouter := router.Group("/post")
	{
		postRouter.POST("/", middleware.Authentication(), controllers.CreatePost)
		postRouter.GET("/", middleware.Authentication(), controllers.GetPosts)

		// Comment post
		postRouter.POST("/:id/comment", middleware.Authentication(), controllers.CreateComment)

		// Like Post
		postRouter.POST("/:id/like", middleware.Authentication(), controllers.CreateLikePost)
		postRouter.POST("/:id/comment/like", middleware.Authentication(), controllers.CreateLikeComment)
	}
	return router
}
