package handlers

import (
	"b2match/backend/dto"
	"b2match/backend/models"

	"github.com/gin-gonic/gin"
)

// GET /companies
func GetCompanies(c *gin.Context) {
	getResources[models.Company](c)
}

// GET /companies/:id
func GetCompanyByID(c *gin.Context) {
	getResourceByID[models.Company](c)
}

// POST /companies
func CreateCompany(c *gin.Context) {
	var input dto.NewCompanyJSON
	if err := bindJSON(c, &input); err != nil {
		return
	}

	company := models.Company{
		Name:     input.Name,
		Location: input.Location,
		About:    input.About,
	}

	createResource(c, &company)
}

// PATCH /companies/:id
func UpdateCompany(c *gin.Context) {
	var company models.Company
	if err := findResourceByID(c, &company, c.Param("id")); err != nil {
		return
	}

	var input dto.UpdateCompanyJSON
	if err := bindJSON(c, &input); err != nil {
		return
	}

	updateResource(c, &company, &input)
}

// DELETE /companies/:id
func DeleteCompany(c *gin.Context) {
	deleteResource[models.Company](c, "Representatives")
}

// GET /company/:id/representatives
func GetCompanyRepresentatives(c *gin.Context) {
	getNestedResources[models.Company, models.User](c, "Representatives")
}
