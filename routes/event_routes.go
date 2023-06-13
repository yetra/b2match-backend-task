package routes

import (
	"b2match/backend/handlers"

	"github.com/gin-gonic/gin"
)

func addEventGroup(router *gin.Engine) {
	eventRouter := router.Group("/events")
	{
		eventRouter.GET("", handlers.GetEvents)
		eventRouter.GET("/:id", handlers.GetEventByID)
		eventRouter.POST("", handlers.CreateEvent)
		eventRouter.PATCH("/:id", handlers.UpdateEvent)
		eventRouter.DELETE("/:id", handlers.DeleteEvent)

		eventRouter.POST("/:id/join", handlers.JoinEvent)

		eventRouter.GET("/:id/participants", handlers.GetEventParticipants)

		eventRouter.GET("/:id/meetings", handlers.GetEventMeetings)
		eventRouter.POST("/:id/meetings", handlers.CreateEventMeeting)
	}
}
