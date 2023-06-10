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

type joinEventJSON struct {
	ID uint `binding:"required" json:"id"`
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
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"event": event})
}

// GET /events/:id/participants
func FindEventParticipants(c *gin.Context) {
	var event models.Event
	var participants []models.User

	id := c.Param("id")

	if err := database.DB.First(&event, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	err := database.DB.Model(&event).Association("Participants").Find(&participants)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"participants": participants})
}

// POST /events/:id/join
func JoinEvent(c *gin.Context) {
	var event models.Event
	var joinData joinEventJSON
	var participant models.User

	id := c.Param("id")

	if err := database.DB.First(&event, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if err := c.ShouldBindJSON(&joinData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.First(&participant, joinData.ID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	err := database.DB.Model(&event).Association("Participants").Append(&participant)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusCreated)
}
