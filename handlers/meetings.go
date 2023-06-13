package handlers

import (
	"net/http"

	"b2match/backend/database"
	"b2match/backend/dto"
	"b2match/backend/models"

	"github.com/gin-gonic/gin"
)

// GetMeetingByID godoc
//
// @Summary		 Get a single meeting by id
// @Description	 Returns the meeting whose ID value matches the id parameter.
// @Tags		 meetings
// @Produce		 json
// @Param		 id path int true "Meeting ID"
// @Success		 200	{object}	models.Meeting
// @Failure		 404	{object}	dto.Error
// @Router		 /meetings/{id} [get]
func GetMeetingByID(c *gin.Context) {
	getResourceByID[models.Meeting](c)
}

// ScheduleMeeting godoc
//
// @Summary		 Schedule a meeting
// @Description	 Marks a meeting as scheduled if all its invites are accepted. Returns the scheduled meeting.
// @Tags		 meetings
// @Produce		 json
// @Param		 id path int true "Meeting ID"
// @Success		 200	{object}	models.Meeting
// @Failure		 400	{object}	dto.Error
// @Failure		 404	{object}	dto.Error
// @Failure		 422	{object}	dto.Error
// @Router		 /meetings/{id}/schedule [patch]
func ScheduleMeeting(c *gin.Context) {
	var meeting models.Meeting
	if err := findResourceByID(c, &meeting, c.Param("id")); err != nil {
		return
	}

	var invites []models.Invite
	if err := findNestedResources(c, &meeting, &invites, "Invites"); err != nil {
		return
	}

	for _, invite := range invites {
		if invite.Status != models.Accepted {
			err_message := "found an invite of status Pending or Rejected"
			c.JSON(http.StatusUnprocessableEntity, dto.Error{Errors: err_message})
			return
		}
	}

	err := database.DB.Model(&meeting).Update("Scheduled", true).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Error{Errors: err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

// DeleteMeeting godoc
//
// @Summary      Deletes a meeting
// @Description  Deletes a meeting and its invites.
// @Tags         meetings
// @Accept       json
// @Produce      json
// @Param		 id	path int true "Meeting ID"
// @Success      204  {object}  nil
// @Failure      404  {object}  dto.Error
// @Router       /meetings/{id} [delete]
func DeleteMeeting(c *gin.Context) {
	deleteResource[models.Meeting](c, "Invites")
}

// GetMeetingInvites godoc
//
// @Summary      Get meeting invites
// @Description  Responds with a list of the meeting's invites as JSON.
// @Tags         meetings
// @Produce      json
// @Param		 id	path int true "Meeting ID"
// @Success      200	{array}		models.Invite
// @Failure      404	{object}	dto.Error
// @Failure      500	{object}	dto.Error
// @Router       /meetings/{id}/invites [get]
func GetMeetingInvites(c *gin.Context) {
	getNestedResources[models.Meeting, models.Invite](c, "Invites")
}

// CreateMeetingInvite godoc
//
// @Summary      Create a new meeting invite
// @Description  Creates an invite for the meeting specified by id and stores it in the database. Returns the new invite.
// @Tags         meetings
// @Accept       json
// @Produce      json
// @Success      201 	{object}	models.Invite
// @Failure      400 	{object}	dto.Error
// @Failure      500 	{object}	dto.Error
// @Router       /meeting/{id}/invites [post]
func CreateMeetingInvite(c *gin.Context) {
	var input dto.NewInviteJSON
	if err := bindJSON(c, &input); err != nil {
		return
	}

	var meeting models.Meeting
	if err := findResourceByID(c, &meeting, c.Param("id")); err != nil {
		return
	}

	var invitee models.User
	if err := findResourceByID(c, &invitee, input.InviteeID); err != nil {
		return
	}

	if err := checkInviteeIsNotOrganizer(invitee.ID, meeting.OrganizerID); err != nil {
		c.JSON(http.StatusUnprocessableEntity, dto.Error{Errors: err.Error()})
		return
	}
	if err := checkInviteeIsAParticipant(invitee.ID, meeting.EventID); err != nil {
		c.JSON(http.StatusUnprocessableEntity, dto.Error{Errors: err.Error()})
		return
	}
	if err := checkInviteeNotAlreadyInvited(invitee.ID, meeting.ID); err != nil {
		c.JSON(http.StatusUnprocessableEntity, dto.Error{Errors: err.Error()})
		return
	}

	invite := models.Invite{
		Status: models.Pending,

		MeetingID: meeting.ID,
		InviteeID: invitee.ID,
	}

	createResource(c, &invite)
}
