package main

import (
	"b2match/backend/database"
	"b2match/backend/routes"
)

func main() {
	database.SetUpDB("b2match.db")

	router := routes.CreateRouter()
	router.Run(":8085")
}
