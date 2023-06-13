package handlers

import (
	"errors"
	"net/http"

	"b2match/backend/database"
	"b2match/backend/dto"
	"b2match/backend/models"

	"github.com/gin-gonic/gin"
)

// GET /event/:id/meetings
func GetEventMeetings(c *gin.Context) {
	getNestedResources[models.Event, models.Meeting](c, "Meetings")
}

// POST /events/:id/meetings
func CreateEventMeeting(c *gin.Context) {
	var newMeeting dto.NewMeetingJSON
	if err := bindJSON(c, &newMeeting); err != nil {
		return
	}

	var event models.Event
	if err := findResourceByID(c, &event, c.Param("id")); err != nil {
		return
	}

	if err := checkNewMeetingIsDuringEvent(newMeeting, event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
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
			err_message := "found an invite of status Pending or Rejected"
			c.JSON(http.StatusUnprocessableEntity, gin.H{"errors": err_message})
			return
		}
	}

	err := database.DB.Model(&meeting).Update("Scheduled", true).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

// DELETE /meetings/:id
func DeleteMeeting(c *gin.Context) {
	deleteResource[models.Meeting](c, "Invites")
}

func checkNewMeetingIsDuringEvent(meeting dto.NewMeetingJSON, event models.Event) error {
	if meeting.StartTime.After(event.StartDate) && meeting.EndTime.Before(event.EndDate) {
		return nil
	}

	return errors.New("meeting must happen during the event")
}
