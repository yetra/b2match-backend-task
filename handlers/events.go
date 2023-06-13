package handlers

import (
	"net/http"

	"b2match/backend/database"
	"b2match/backend/dto"
	"b2match/backend/models"

	"github.com/gin-gonic/gin"
)

// GET /events
func GetEvents(c *gin.Context) {
	getResources[models.Event](c)
}

// GET /events/:id
func GetEventByID(c *gin.Context) {
	getResourceByID[models.Event](c)
}

// POST /events
func CreateEvent(c *gin.Context) {
	var newEvent dto.NewEventJSON
	if err := bindJSON(c, &newEvent); err != nil {
		return
	}

	event := models.Event{
		Name:      newEvent.Name,
		Location:  newEvent.Location,
		Agenda:    newEvent.Agenda,
		StartDate: newEvent.StartDate,
		EndDate:   newEvent.EndDate,
	}

	createResource(c, &event)
}

// PATCH /events/:id
func UpdateEvent(c *gin.Context) {
	var event models.Event
	if err := findResourceByID(c, &event, c.Param("id")); err != nil {
		return
	}

	var updatedEvent dto.UpdateEventJSON
	if err := bindJSON(c, &updatedEvent); err != nil {
		return
	}

	updateResource(c, &event, &updatedEvent)
}

// DELETE /events/:id
func DeleteEvent(c *gin.Context) {
	deleteResource[models.Event](c, "Meetings")
}

// POST /events/:id/join
func JoinEvent(c *gin.Context) {
	var event models.Event
	if err := findResourceByID(c, &event, c.Param("id")); err != nil {
		return
	}

	var joinData dto.JoinEventJSON
	if err := bindJSON(c, &joinData); err != nil {
		return
	}

	var participant models.User
	if err := findResourceByID(c, &participant, joinData.ID); err != nil {
		return
	}

	err := database.DB.Model(&event).Association("Participants").Append(&participant)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": err.Error()})
		return
	}

	c.Status(http.StatusCreated)
}

// GET /events/:id/participants
func GetEventParticipants(c *gin.Context) {
	getNestedResources[models.Event, models.User](c, "Participants")
}
