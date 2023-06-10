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
	Location string
	Agenda   string

	StartDate time.Time `binding:"required,ltefield=EndDate" json:"start_date"`
	EndDate   time.Time `binding:"required" json:"end_date"`
}

type joinEventJSON struct {
	ID uint `binding:"required" json:"id"`
}

// GET /events
func FindEvents(c *gin.Context) {
	var events []models.Event
	database.DB.Preload("Participants").Find(&events)

	c.JSON(http.StatusOK, gin.H{"events": events})
}

// GET /events/:id
func FindEventByID(c *gin.Context) {
	var event models.Event

	id := c.Param("id")

	if err := database.DB.Preload("Participants").First(&event, id).Error; err != nil {
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
