package handlers

import (
	"net/http"

	"b2match/backend/database"
	"b2match/backend/dto"
	"b2match/backend/models"

	"github.com/gin-gonic/gin"
)

// GetEvents godoc
//
// @Summary		 Get events
// @Description	 Responds with a list of all events as JSON.
// @Tags		 events
// @Produce		 json
// @Success		 200	{array}		models.Event
// @Router		 /events [get]
func GetEvents(c *gin.Context) {
	getResources[models.Event](c)
}

// GetEventByID godoc
//
// @Summary		 Get a single event by id
// @Description	 Returns the event whose ID value matches the id parameter.
// @Tags		 events
// @Produce		 json
// @Param		 id path int true "Event ID"
// @Success		 200	{object}	models.Event
// @Failure		 404	{object}	gin.H
// @Router		 /events/{id} [get]
func GetEventByID(c *gin.Context) {
	getResourceByID[models.Event](c)
}

// CreateEvent godoc
//
// @Summary      Create a new event
// @Description  Creates an event from the input JSON and stores it in the database. Returns the new event.
// @Tags         events
// @Accept       json
// @Produce      json
// @Success      201 	{object}	models.Event
// @Failure      400 	{object}	gin.H
// @Failure      500 	{object}	gin.H
// @Router       /events [post]
func CreateEvent(c *gin.Context) {
	var input dto.NewEventJSON
	if err := bindJSON(c, &input); err != nil {
		return
	}

	event := models.Event{
		Name:      input.Name,
		Location:  input.Location,
		Agenda:    input.Agenda,
		StartDate: input.StartDate,
		EndDate:   input.EndDate,
	}

	createResource(c, &event)
}

// UpdateEvent godoc
//
// @Summary      Update an existing event
// @Description  Updates an event with the input JSON. Returns the updated event.
// @Tags         events
// @Accept       json
// @Produce      json
// @Param		 id	path int true "Event ID"
// @Success      200	{object}	models.Event
// @Failure      400	{object}	gin.H
// @Failure      404	{object}	gin.H
// @Failure      500	{object}	gin.H
// @Router       /events/{id} [patch]
func UpdateEvent(c *gin.Context) {
	var event models.Event
	if err := findResourceByID(c, &event, c.Param("id")); err != nil {
		return
	}

	var input dto.UpdateEventJSON
	if err := bindJSON(c, &input); err != nil {
		return
	}

	updateResource(c, &event, &input)
}

// DeleteEvent godoc
//
// @Summary      Delete an event
// @Description  Deletes an event and its meetings.
// @Tags         events
// @Accept       json
// @Produce      json
// @Param		 id	path int true "Event ID"
// @Success      204  {object}  nil
// @Failure      404  {object}  gin.H
// @Router       /events/{id} [delete]
func DeleteEvent(c *gin.Context) {
	deleteResource[models.Event](c, "Meetings")
}

// JoinEvent godoc
//
// @Summary      Join an event
// @Description  Adds the user specified in the request JSON to the event's participants.
// @Tags         events
// @Accept       json
// @Produce      json
// @Param		 id	path int true "Event ID"
// @Success      201  {object}  nil
// @Failure      400  {object}  gin.H
// @Failure      404  {object}  gin.H
// @Router       /events/{id}/join [post]
func JoinEvent(c *gin.Context) {
	var event models.Event
	if err := findResourceByID(c, &event, c.Param("id")); err != nil {
		return
	}

	var input dto.JoinEventJSON
	if err := bindJSON(c, &input); err != nil {
		return
	}

	var participant models.User
	if err := findResourceByID(c, &participant, input.ID); err != nil {
		return
	}

	err := database.DB.Model(&event).Association("Participants").Append(&participant)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": err.Error()})
		return
	}

	c.Status(http.StatusCreated)
}

// GetEventParticipants godoc
//
// @Summary      Get event participants
// @Description  Responds with a list of event participants as JSON.
// @Tags         events
// @Produce      json
// @Param		 id	path int true "Event ID"
// @Success      200	{array}		models.User
// @Failure      404	{object}	gin.H
// @Failure      500	{object}	gin.H
// @Router       /events/{id}/participants [get]
func GetEventParticipants(c *gin.Context) {
	getNestedResources[models.Event, models.User](c, "Participants")
}
