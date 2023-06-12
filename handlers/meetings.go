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
func GetEventMeetings(c *gin.Context) {
	getNestedResources[models.Event, models.Meeting](c, "Meetings")
}

// POST /events/:id/meetings
func CreateEventMeeting(c *gin.Context) {
	var newMeeting newMeetingJSON

	if err := bindJSON(c, &newMeeting); err != nil {
		return
	}

	var event models.Event

	if err := findResourceByID(c, &event, c.Param("id")); err != nil {
		return
	}

	var organizer models.User

	if err := findResourceByID(c, &organizer, newMeeting.OrganizerID); err != nil {
		return
	}

	meeting := models.Meeting{
		StartTime:   newMeeting.StartTime,
		EndTime:     newMeeting.EndTime,
		EventID:     event.ID,
		OrganizerID: organizer.ID,
	}

	createResource(c, &meeting)
}

// GET /meetings/:id
func GetMeetingByID(c *gin.Context) {
	getResourceByID[models.Meeting](c)
}

// PATCH /meetings/:id/schedule
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
			err_message := "Found an invite of status Pending or Rejected."
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err_message})
			return
		}
	}

	err := database.DB.Model(&meeting).Update("Scheduled", true).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}
