package main

import (
	"os"
	"twittir-go/database"
	"twittir-go/router"

	"github.com/gin-gonic/gin"
)

func main() {
	database.StartDB()
	gin.SetMode(gin.ReleaseMode)

	var PORT = os.Getenv("PORT")

	router.StartServer().Run("0.0.0.0:" + PORT)
}
