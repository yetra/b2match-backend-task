package handlers

import (
	"b2match/backend/models"

	"github.com/gin-gonic/gin"
)

type newInviteJSON struct {
	InviteeID uint `binding:"required" json:"invitee_id"`
}

type rsvpJSON struct {
	Status models.Status `binding:"required,min=1,max=2" json:"status"`
}

// GET /meetings/:id/invites
func GetMeetingInvites(c *gin.Context) {
	getNestedResources[models.Meeting, models.Invite](c, "Invites")
}

// POST /meetings/:id/invites
func CreateMeetingInvite(c *gin.Context) {
	var inviteData newInviteJSON

	if err := bindJSON(c, &inviteData); err != nil {
		return
	}

	var meeting models.Meeting

	if err := findResourceByID(c, &meeting, c.Param("id")); err != nil {
		return
	}

	var invitee models.User

	if err := findResourceByID(c, &invitee, inviteData.InviteeID); err != nil {
		return
	}

	invite := models.Invite{
		Status: models.Pending,

		MeetingID: meeting.ID,
		InviteeID: invitee.ID,
	}

	createResource(c, &invite)
}

// GET /invites/:id
func GetInviteByID(c *gin.Context) {
	getResourceByID[models.Meeting](c)
}

// PATCH /invites/:id/rsvp
func RespondToInvite(c *gin.Context) {
	var rsvpData rsvpJSON

	if err := bindJSON(c, &rsvpData); err != nil {
		return
	}

	var invite models.Invite

	if err := findResourceByID(c, &invite, c.Param("id")); err != nil {
		return
	}

	updateResource(c, &invite, &rsvpData)
}
