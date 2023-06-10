package handlers

import (
	"net/http"

	"b2match/backend/database"
	"b2match/backend/models"

	"github.com/gin-gonic/gin"
)

// GET /companies
func FindCompanies(c *gin.Context) {
	var companies []models.Company
	database.DB.Find(&companies)

	c.JSON(http.StatusOK, gin.H{"companies": companies})
}

// GET /companies/:id
func FindCompanyByID(c *gin.Context) {
	var company models.Company

	id := c.Param("id")

	if err := database.DB.First(&company, id).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"company": company})
}

// POST /companies
func CreateCompany(c *gin.Context) {
	var company models.Company

	if err := c.ShouldBindJSON(&company); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.Create(&company).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"company": company})
}
