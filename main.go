package main

import (
	"b2match/backend/database"
	"b2match/backend/routes"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	database.SetUpDB("b2match.db", &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	router := routes.CreateRouter()
	router.Run(":8085")
}
