package handlers

import (
	"net/http"
	"time"

	"b2match/backend/database"
	"b2match/backend/models"

	"github.com/gin-gonic/gin"
)

type newEventJSON struct {
	Name     string `binding:"required"`
	Location string `binding:"required"`
	Agenda   string

	StartDate time.Time `binding:"required,ltefield=EndDate" time_format:"2006-01-02"`
	EndDate   time.Time `binding:"required" time_format:"2006-01-02"`
}

// GET /events
func FindEvents(c *gin.Context) {
	var events []models.Event
	database.DB.Find(&events)

	c.JSON(http.StatusOK, gin.H{"events": events})
}

// GET /events/:id
func FindEventByID(c *gin.Context) {
	var event models.Event

	id := c.Param("id")

	if err := database.DB.First(&event, id).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"event": event})
}

// POST /events
func CreateEvent(c *gin.Context) {
	var newEvent newEventJSON

	if err := c.ShouldBindJSON(&newEvent); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	event := models.Event{
		Name:      newEvent.Name,
		Location:  newEvent.Location,
		Agenda:    newEvent.Agenda,
		StartDate: newEvent.StartDate,
		EndDate:   newEvent.EndDate,
	}
	if err := database.DB.Create(&event).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"event": event})
}
