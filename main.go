package main

import (
	"b2match/backend/database"
	"b2match/backend/docs"
	"b2match/backend/routes"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setUpSwagger(router *gin.Engine) {
	docs.SwaggerInfo.Title = "b2match API"
	docs.SwaggerInfo.Description = "An API for event management."
	docs.SwaggerInfo.Version = "1.0"

	docs.SwaggerInfo.Host = "localhost:8085"
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Schemes = []string{"http"}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

func main() {
	database.SetUpDB("b2match.db", &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	router := routes.CreateRouter()
	setUpSwagger(router)

	router.Run(":8085")
}
