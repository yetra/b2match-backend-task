package handlers

import (
	"net/http"
	"time"

	"b2match/backend/database"
	"b2match/backend/models"

	"github.com/gin-gonic/gin"
)

type newMeetingJSON struct {
	StartTime time.Time `binding:"required,ltefield=EndTime" json:"start_time"`
	EndTime   time.Time `binding:"required" json:"end_time"`

	OrganizerID uint `binding:"required" json:"organizer_id"`
}

// GET /event/:id/meetings
func FindEventMeetings(c *gin.Context) {
	findNestedResources[models.Event, models.Meeting](c, "Meetings")
}

// GET /event/:id/meetings/:meeting_id
func FindEventMeetingByID(c *gin.Context) {
	var meeting models.Meeting

	id := c.Param("id")
	meeting_id := c.Param("meeting_id")

	err := database.DB.Where("event_id = ?", id).First(&meeting, meeting_id).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"meeting": meeting})
}

// POST /events/:id/meetings
func CreateEventMeeting(c *gin.Context) {
	var newMeeting newMeetingJSON

	if err := c.ShouldBindJSON(&newMeeting); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var event models.Event

	id := c.Param("id")

	if err := database.DB.First(&event, id).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var organizer models.User

	if err := database.DB.First(&organizer, newMeeting.OrganizerID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	meeting := models.Meeting{
		StartTime:   newMeeting.StartTime,
		EndTime:     newMeeting.EndTime,
		EventID:     event.ID,
		OrganizerID: organizer.ID,
	}

	if err := database.DB.Create(&meeting).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"meeting": meeting})
}

// PATCH /meetings/:id/schedule
func ScheduleMeeting(c *gin.Context) {
	var meeting models.Meeting

	id := c.Param("id")

	if err := database.DB.First(&meeting, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	var invites []models.Invite

	err := database.DB.Model(&meeting).Association("Invites").Find(&invites)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for _, invite := range invites {
		if invite.Status != models.Accepted {
			err_message := "Found an invite of status Pending or Rejected."
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err_message})
			return
		}
	}

	err = database.DB.Model(&meeting).Update("Scheduled", true).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}
