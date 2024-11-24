package main

import (
	"twittir-go/config"
	"twittir-go/database"
	"twittir-go/router"
)

const PORT = ":3001"

func main() {
	database.StartDB()

	config.GoogleConfig()

	router.StartServer().Run(PORT)
}
