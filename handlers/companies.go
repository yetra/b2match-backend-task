package handlers

import (
	"b2match/backend/models"

	"github.com/gin-gonic/gin"
)

type newCompanyJSON struct {
	Name     string `binding:"required"`
	Location string `binding:"required"`
	About    string
}

type updateCompanyJSON struct {
	Location string
	About    string
}

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
	var newCompany newCompanyJSON
	if err := bindJSON(c, &newCompany); err != nil {
		return
	}

	company := models.Company{
		Name:     newCompany.Name,
		Location: newCompany.Location,
		About:    newCompany.About,
	}

	createResource(c, &company)
}

// PATCH /companies/:id
func UpdateCompany(c *gin.Context) {
	var company models.Company
	if err := findResourceByID(c, &company, c.Param("id")); err != nil {
		return
	}

	var updatedCompany updateCompanyJSON
	if err := bindJSON(c, &updatedCompany); err != nil {
		return
	}

	updateResource(c, &company, &updatedCompany)
}

// DELETE /companies/:id
func DeleteCompany(c *gin.Context) {
	deleteResource[models.Company](c, "Representatives")
}

// GET /company/:id/representatives
func GetCompanyRepresentatives(c *gin.Context) {
	getNestedResources[models.Company, models.User](c, "Representatives")
}
