// @title Twittir API
// @version 1.0
// @description This is a sample API for Twittir.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@swagger.io

// @host localhost:3001
// @BasePath /

package routes

import (
	"twittir-go/internal/database"
	"twittir-go/internal/handler"
	"twittir-go/internal/repositories"
	"twittir-go/internal/services"

	// controllers "twittir-go/internal/handler"
	"twittir-go/internal/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	_ "twittir-go/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func StartServer() *gin.Engine {

	config := cors.DefaultConfig()
	// config.AllowOrigins = []string{"http://localhost:3000"} // Replace with your allowed origins
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}

	router := gin.Default()
	router.Use(cors.New(config))

	// Database, Repositories, Services and Handlers Call
	db := database.GetDB()
	userRepo := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	postRepo := repositories.NewPostRepository(db)
	postService := services.NewPostService(postRepo)
	postHandler := handler.NewPostHandler(postService, userHandler)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	apiRouter := router.Group("/api")
	{
		// Auth
		apiRouter.POST("/register", userHandler.Register)
		apiRouter.POST("/login", userHandler.Login)
		apiRouter.GET("/users/profile", middleware.Authentication(), userHandler.ShowProfile)
		apiRouter.PUT("/users/update", middleware.Authentication(), userHandler.UpdateProfile)
		apiRouter.GET("/users/:id", middleware.Authentication(), userHandler.GetUserByID)
		apiRouter.GET("/search", middleware.Authentication(), userHandler.SearchUser)

		// // Post
		apiRouter.POST("/posts", middleware.Authentication(), postHandler.CreatePost)
		apiRouter.GET("/posts", middleware.Authentication(), postHandler.GetPosts)
		apiRouter.GET("/posts/:id", middleware.Authentication(), postHandler.GetPostByID)
		apiRouter.GET("/posts/user/:id", middleware.Authentication(), postHandler.GetPostByUserID)
		// apiRouter.DELETE("/posts/:id", middleware.Authentication(), controllers.DeletePost)

		// // Comment post
		// apiRouter.POST("/posts/:id/comment", middleware.Authentication(), controllers.CreateComment)
		// apiRouter.DELETE("/comments/:id", middleware.Authentication(), controllers.DeleteComment)

		// // Like Post
		// apiRouter.POST("/posts/:id/like", middleware.Authentication(), controllers.CreateLikePost)
		// apiRouter.POST("/comments/:id/like", middleware.Authentication(), controllers.CreateLikeComment)

		// // Relationship / Follow
		// apiRouter.POST("/users/:id/follow", middleware.Authentication(), controllers.FollowUser)
	}

	return router
}
