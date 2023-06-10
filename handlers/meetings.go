package handlers

import (
	"net/http"
	"time"

	"b2match/backend/database"
	"b2match/backend/models"

	"github.com/gin-gonic/gin"
)

type newMeetingJSON struct {
	StartTime time.Time `binding:"required,ltefield=EndTime" json:"start_time" time_format:"2006-01-02"`
	EndTime   time.Time `binding:"required" json:"end_time" time_format:"2006-01-02"`

	OrganizerID uint `binding:"required" json:"organizer_id"`
}

// GET /event/:id/meetings
func FindEventMeetings(c *gin.Context) {
	var event models.Event
	var meetings []models.Meeting

	id := c.Param("id")

	if err := database.DB.First(&event, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	err := database.DB.Model(&event).Association("Meetings").Find(&meetings)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"meetings": meetings})
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
