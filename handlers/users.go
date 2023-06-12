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

type updateUserJSON struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`

	Location string
	About    string

	Password string
}

// GET /users
func GetUsers(c *gin.Context) {
	getResources[models.User](c)
}

// GET /users/:id
func GetUserByID(c *gin.Context) {
	getResourceByID[models.User](c)
}

// POST /users
func CreateUser(c *gin.Context) {
	var newUser newUserJSON

	if err := bindJSON(c, &newUser); err != nil {
		return
	}

	var company models.Company

	if err := findResourceByID(c, &company, newUser.CompanyID); err != nil {
		return
	}

	user := models.User{
		FirstName: newUser.FirstName,
		LastName:  newUser.LastName,
		Location:  newUser.Location,
		About:     newUser.About,
		EMail:     newUser.EMail,
		Password:  newUser.Password,
		CompanyID: company.ID,
	}

	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"user": user})
}

// GET /users/:id/meetings
func GetUserScheduledMeetings(c *gin.Context) {
	var user models.User

	if err := findResourceByID(c, &user, c.Param("id")); err != nil {
		return
	}

	var invites []models.Invite

	err := database.DB.Model(&user).Association("Invites").Find(&invites)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var scheduledMeetings []models.Meeting

	for _, invite := range invites {
		var meeting models.Meeting

		if err := findResourceByID(c, &invite, invite.MeetingID); err != nil {
			return
		}

		if meeting.Scheduled {
			scheduledMeetings = append(scheduledMeetings, meeting)
		}
	}

	c.JSON(http.StatusOK, gin.H{"meetings": scheduledMeetings})
}

// GET /users/:id/invites
func GetUserInvites(c *gin.Context) {
	getNestedResources[models.User, models.Invite](c, "Invites")
}

// PATCH /users/:id
func UpdateUser(c *gin.Context) {
	var user models.User

	if err := findResourceByID(c, &user, c.Param("id")); err != nil {
		return
	}

	var updatedUser updateUserJSON

	if err := bindJSON(c, &updatedUser); err != nil {
		return
	}

	if err := database.DB.Model(&user).Updates(&updatedUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}
