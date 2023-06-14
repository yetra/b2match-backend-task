package handlers

import (
	"net/http"

	"b2match/backend/dto"
	"b2match/backend/models"

	"github.com/gin-gonic/gin"
)

// GetUsers godoc
//
// @Summary		 Get users
// @Description	 Responds with a list of all users as JSON.
// @Tags		 users
// @Produce		 json
// @Success		 200	{array}		models.User
// @Router		 /users [get]
func GetUsers(c *gin.Context) {
	getResources[models.User](c)
}

// GetUserByID godoc
//
// @Summary		 Get a single user by id
// @Description	 Returns the user whose ID value matches the id parameter.
// @Tags		 users
// @Produce		 json
// @Param		 id path int true "User ID"
// @Success		 200	{object}	models.User
// @Failure		 404	{object}	dto.Error
// @Router		 /users/{id} [get]
func GetUserByID(c *gin.Context) {
	getResourceByID[models.User](c)
}

// CreateUser godoc
//
// @Summary      Create a new user
// @Description  Creates a user from the input JSON and stores it in the database. Returns the new user.
// @Tags         users
// @Accept       json
// @Produce      json
// @Success      201 	{object}	models.User
// @Failure      400 	{object}	dto.Error
// @Failure      500 	{object}	dto.Error
// @Router       /users [post]
func CreateUser(c *gin.Context) {
	input, err := bindJSON[dto.NewUserJSON](c)
	if err != nil {
		return
	}

	company, err := findResourceByID[models.Company](c, input.CompanyID)
	if err != nil {
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

// UpdateUser godoc
//
// @Summary      Update an existing user
// @Description  Updates a user with the input JSON. Returns the updated user.
// @Tags         users
// @Accept       json
// @Produce      json
// @Param		 id	path int true "User ID"
// @Success      200	{object}	models.User
// @Failure      400	{object}	dto.Error
// @Failure      404	{object}	dto.Error
// @Failure      500	{object}	dto.Error
// @Router       /users/{id} [patch]
func UpdateUser(c *gin.Context) {
	user, err := findResourceByID[models.User](c, c.Param("id"))
	if err != nil {
		return
	}

	input, err := bindJSON[dto.UpdateUserJSON](c)
	if err != nil {
		return
	}

	updateResource(c, &user, &input)
}

// DeleteUser godoc
//
// @Summary      Delete a user
// @Description  Deletes a user, its organized meetings, and invites.
// @Tags         users
// @Accept       json
// @Produce      json
// @Param		 id	path int true "User ID"
// @Success      204  {object}  nil
// @Failure      404  {object}  dto.Error
// @Router       /users/{id} [delete]
func DeleteUser(c *gin.Context) {
	deleteResource[models.User](c, []string{"OrganizedMeetings", "Invites"})
}

// GetUserScheduledMeetings godoc
//
// @Summary      Get user scheduled meetings
// @Description  Responds with a list of the user's scheduled meetings as JSON.
// @Tags         users
// @Produce      json
// @Param		 id	path int true "User ID"
// @Success      200	{array}		models.Meeting
// @Failure      404	{object}	dto.Error
// @Failure      500	{object}	dto.Error
// @Router       /users/{id}/scheduled-meetings [get]
func GetUserScheduledMeetings(c *gin.Context) {
	user, err := findResourceByID[models.User](c, c.Param("id"))
	if err != nil {
		return
	}

	invites, err := findNestedResources[models.Invite](c, &user, "Invites")
	if err != nil {
		return
	}

	var scheduledMeetings []models.Meeting
	for _, invite := range invites {

		meeting, err := findResourceByID[models.Meeting](c, invite.MeetingID)
		if err != nil {
			return
		}

		if meeting.Scheduled {
			scheduledMeetings = append(scheduledMeetings, meeting)
		}
	}

	c.JSON(http.StatusOK, scheduledMeetings)
}

// GetUserInvites godoc
//
// @Summary      Get user invites
// @Description  Responds with a list of the user's meeting invites as JSON.
// @Tags         users
// @Produce      json
// @Param		 id	path int true "User ID"
// @Success      200	{array}		models.Invite
// @Failure      404	{object}	dto.Error
// @Failure      500	{object}	dto.Error
// @Router       /users/{id}/invites [get]
func GetUserInvites(c *gin.Context) {
	getNestedResources[models.User, models.Invite](c, "Invites")
}
