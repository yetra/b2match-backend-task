package handlers

import (
	"net/http"

	"b2match/backend/database"
	"b2match/backend/models"

	"github.com/gin-gonic/gin"
)

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
