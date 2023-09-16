package main

import (
	"os"
	"twittir-go/database"
	"twittir-go/router"
)

func main() {
	database.StartDB()

	var PORT = os.Getenv("PORT")

	router.StartServer().Run(PORT)
}
