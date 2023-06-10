package main

import (
	"b2match/backend/database"
	"b2match/backend/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	database.SetUpDB()

	route := gin.Default()

	route.GET("/companies", handlers.FindCompanies)
	route.GET("/companies/:id", handlers.FindCompanyByID)
	route.POST("/companies", handlers.CreateCompany)

	route.Run(":8085")
}
