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

	route.GET("/users", handlers.FindUsers)
	route.GET("/users/:id", handlers.FindUserByID)
	route.POST("/users", handlers.CreateUser)

	route.GET("/events", handlers.FindEvents)
	route.GET("/events/:id", handlers.FindEventByID)
	route.POST("/events", handlers.CreateUser)

	route.Run(":8085")
}
