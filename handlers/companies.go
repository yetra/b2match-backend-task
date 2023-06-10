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
