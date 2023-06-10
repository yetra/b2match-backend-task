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

// GET /event/:id/meetings/:meeting_id/invites
func FindMeetingInvites(c *gin.Context) {
	var event models.Event
	var meeting models.Meeting
	var invites []models.Invite

	id := c.Param("id")

	if err := database.DB.First(&event, id).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	meeting_id := c.Param("meeting_id")

	if err := database.DB.First(&meeting, meeting_id).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := database.DB.Model(&meeting).Association("Invites").Find(&invites)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"invites": invites})
}

// POST /events/:id/meetings/:meeting_id/invites
func CreateMeetingInvite(c *gin.Context) {
	var inviteData newInviteJSON

	if err := c.ShouldBindJSON(&inviteData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var event models.Event

	id := c.Param("id")

	if err := database.DB.First(&event, id).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var meeting models.Meeting

	meeting_id := c.Param("meeting_id")

	if err := database.DB.First(&meeting, meeting_id).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var invitee models.User

	if err := database.DB.First(&invitee, inviteData.InviteeID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
