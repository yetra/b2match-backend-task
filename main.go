package main

import (
	"b2match/backend/database"
	"b2match/backend/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	database.SetUpDB("b2match.db")

	router := gin.Default()

	routes.AddRoutes(router)

	router.Run(":8085")
}
