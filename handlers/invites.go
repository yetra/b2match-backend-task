package handlers

import (
	"net/http"

	"b2match/backend/database"
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

// GET /meetings/:id/invites/:invite_id
func GetMeetingInviteByID(c *gin.Context) {
	var invite models.Invite

	id := c.Param("id")
	invite_id := c.Param("invite_id")

	err := database.DB.Where("meeting_id = ?", id).First(&invite, invite_id).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"invite": invite})
}

// POST /meetings/:id/invites
func CreateMeetingInvite(c *gin.Context) {
	var inviteData newInviteJSON

	if err := c.ShouldBindJSON(&inviteData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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

	if err := database.DB.Create(&invite).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"invite": invite})
}

// PATCH /invites/:id/rsvp
func RespondToInvite(c *gin.Context) {
	var rsvpData rsvpJSON

	if err := c.ShouldBindJSON(&rsvpData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var invite models.Invite

	if err := findResourceByID(c, &invite, c.Param("id")); err != nil {
		return
	}

	err := database.DB.Model(&invite).Update("Status", rsvpData.Status).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
