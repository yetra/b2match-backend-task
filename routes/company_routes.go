package routes

import (
	"b2match/backend/handlers"

	"github.com/gin-gonic/gin"
)

func addCompanyGroup(router *gin.Engine) {
	companyRouter := router.Group("/companies")
	{
		companyRouter.GET("", handlers.GetCompanies)
		companyRouter.GET("/:id", handlers.GetCompanyByID)
		companyRouter.POST("", handlers.CreateCompany)
		companyRouter.PATCH("/:id", handlers.UpdateCompany)
		companyRouter.DELETE("/:id", handlers.DeleteCompany)

		companyRouter.GET("/:id/representatives", handlers.GetCompanyRepresentatives)
	}
}
