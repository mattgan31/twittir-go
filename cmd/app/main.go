package main

import (
	"twittir-go/internal/database"
	router "twittir-go/internal/routes"

	_ "twittir-go/docs"
)

const PORT = ":3001"

// @title Twittir Go API
// @version 1.0
// @description Your API description.
// @host localhost:3001
// @BasePath /api/

func main() {
	database.StartDB()

	router.StartServer().Run(PORT)
}
