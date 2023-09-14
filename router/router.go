package router

import (
	"twittir-go/controllers"
	"twittir-go/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func StartServer() *gin.Engine {

	config := cors.DefaultConfig()
	// config.AllowOrigins = []string{"http://localhost:3000"} // Replace with your allowed origins
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}

	router := gin.Default()
	router.Use(cors.New(config))

	apiRouter := router.Group("/api")
	{
		// Auth
		apiRouter.POST("/register", controllers.UserRegister)
		apiRouter.POST("/login", controllers.UserLogin)
		apiRouter.GET("/users/profile", middleware.Authentication(), controllers.GetDetailUser)
		apiRouter.GET("/search", middleware.Authentication(), controllers.SearchUser)

		// Post
		apiRouter.POST("/posts", middleware.Authentication(), controllers.CreatePost)
		apiRouter.GET("/posts", middleware.Authentication(), controllers.GetPosts)
		apiRouter.GET("/posts/:id", middleware.Authentication(), controllers.GetPostByID)
		apiRouter.GET("/posts/user/:id", middleware.Authentication(), controllers.GetPostByUserID)

		// Comment post
		apiRouter.POST("/posts/:id/comment", middleware.Authentication(), controllers.CreateComment)

		// Like Post
		apiRouter.POST("/posts/:id/like", middleware.Authentication(), controllers.CreateLikePost)
		apiRouter.POST("/comments/:id/like", middleware.Authentication(), controllers.CreateLikeComment)

		// Relationship / Follow
		apiRouter.POST("/users/:id/follow", middleware.Authentication(), controllers.FollowUser)
	}

	return router
}
