package handlers

import (
	"net/http"

	"b2match/backend/database"
	"b2match/backend/models"

	"github.com/gin-gonic/gin"
)

type newUserJSON struct {
	FirstName string `binding:"required" json:"first_name"`
	LastName  string `binding:"required" json:"last_name"`

	Location string
	About    string

	EMail    string `binding:"required" json:"e_mail"`
	Password string `binding:"required"`

	CompanyID uint `binding:"required" json:"company_id"`
}

// GET /users
func FindUsers(c *gin.Context) {
	var users []models.User
	database.DB.Find(&users)

	c.JSON(http.StatusOK, gin.H{"users": users})
}

// GET /users/:id
func FindUserByID(c *gin.Context) {
	var user models.User

	id := c.Param("id")

	if err := database.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

// POST /users
func CreateUser(c *gin.Context) {
	var newUser newUserJSON

	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var company models.Company

	if err := database.DB.First(&company, newUser.CompanyID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := models.User{
		FirstName: newUser.FirstName,
		LastName:  newUser.LastName,
		Location:  newUser.Location,
		About:     newUser.About,
		EMail:     newUser.EMail,
		Password:  newUser.Password,
		CompanyID: newUser.CompanyID,
	}

	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"user": user})
}
