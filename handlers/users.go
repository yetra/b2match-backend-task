package handlers

import (
	"net/http"

	"b2match/backend/dto"
	"b2match/backend/models"

	"github.com/gin-gonic/gin"
)

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
	var input dto.NewUserJSON
	if err := bindJSON(c, &input); err != nil {
		return
	}

	var company models.Company
	if err := findResourceByID(c, &company, input.CompanyID); err != nil {
		return
	}

	user := models.User{
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Location:  input.Location,
		About:     input.About,
		EMail:     input.EMail,
		Password:  input.Password,
		CompanyID: company.ID,
	}

	createResource(c, &user)
}

// GET /users/:id/scheduled-meetings
func GetUserScheduledMeetings(c *gin.Context) {
	var user models.User
	if err := findResourceByID(c, &user, c.Param("id")); err != nil {
		return
	}

	var invites []models.Invite
	if err := findNestedResources(c, &user, &invites, "Invites"); err != nil {
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

	c.JSON(http.StatusOK, scheduledMeetings)
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

	var input dto.UpdateUserJSON
	if err := bindJSON(c, &input); err != nil {
		return
	}

	updateResource(c, &user, &input)
}

// DELETE /users/:id
func DeleteUser(c *gin.Context) {
	deleteResource[models.User](c, []string{"OrganizedMeetings", "Invites"})
}
