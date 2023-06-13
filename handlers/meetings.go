package handlers

import (
	"errors"
	"net/http"

	"b2match/backend/database"
	"b2match/backend/dto"
	"b2match/backend/models"

	"github.com/gin-gonic/gin"
)

// GetEventMeetings godoc
//
// @Summary      Get event meetings
// @Description  Responds with a list of the events's meetings as JSON.
// @Tags         events
// @Produce      json
// @Param		 id	path int true "Event ID"
// @Success      200	{array}		models.Meeting
// @Failure      404	{object}	dto.Error
// @Failure      500	{object}	dto.Error
// @Router       /events/{id}/meetings [get]
func GetEventMeetings(c *gin.Context) {
	getNestedResources[models.Event, models.Meeting](c, "Meetings")
}

// CreateEventMeeting godoc
//
// @Summary      Create a new event meetings
// @Description  Creates a meeting for the event specified by id and stores it in the database. Returns the new meeting.
// @Tags         meetings
// @Accept       json
// @Produce      json
// @Success      201 	{object}	models.Meeting
// @Failure      400 	{object}	dto.Error
// @Failure      500 	{object}	dto.Error
// @Router       /events/{id}/meetings [post]
func CreateEventMeeting(c *gin.Context) {
	var input dto.NewMeetingJSON
	if err := bindJSON(c, &input); err != nil {
		return
	}

	var event models.Event
	if err := findResourceByID(c, &event, c.Param("id")); err != nil {
		return
	}

	if err := checkNewMeetingIsDuringEvent(input, event); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error{Errors: err.Error()})
		return
	}

	var organizer models.User
	if err := findResourceByID(c, &organizer, input.OrganizerID); err != nil {
		return
	}

	meeting := models.Meeting{
		StartTime:   input.StartTime,
		EndTime:     input.EndTime,
		EventID:     event.ID,
		OrganizerID: organizer.ID,
	}

	createResource(c, &meeting)
}

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

func checkNewMeetingIsDuringEvent(meeting dto.NewMeetingJSON, event models.Event) error {
	if meeting.StartTime.After(event.StartDate) && meeting.EndTime.Before(event.EndDate) {
		return nil
	}

	return errors.New("meeting must happen during the event")
}
