package handlers

import (
	"net/http"

	"b2match/backend/database"
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

	if err := c.ShouldBindJSON(&newCompany); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	company := models.Company{
		Name:     newCompany.Name,
		Location: newCompany.Location,
		About:    newCompany.About,
	}

	if err := database.DB.Create(&company).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"company": company})
}

// PATCH /companies/:id
func UpdateCompany(c *gin.Context) {
	var company models.Company

	if err := findResourceByID(c, &company, c.Param("id")); err != nil {
		return
	}

	var updatedCompany updateCompanyJSON

	if err := c.ShouldBindJSON(&updatedCompany); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.Model(&company).Updates(&updatedCompany).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"company": company})
}
