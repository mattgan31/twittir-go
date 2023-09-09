package main

import (
	"twittir-go/database"
	"twittir-go/router"
)

const PORT = ":3001"

func main() {
	database.StartDB()

	router.StartServer().Run(PORT)
}
