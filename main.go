package main

import (
	"b2match/backend/database"
	"b2match/backend/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	database.SetUpDB()

	route := gin.Default()

	route.GET("/companies", handlers.GetCompanies)
	route.GET("/companies/:id", handlers.GetCompanyByID)
	route.POST("/companies", handlers.CreateCompany)
	route.PATCH("/companies/:id", handlers.UpdateCompany)

	route.GET("/users", handlers.GetUsers)
	route.GET("/users/:id", handlers.GetUserByID)
	route.POST("/users", handlers.CreateUser)
	route.PATCH("/users/:id", handlers.UpdateUser)

	route.GET("/users/:id/meetings", handlers.GetUserScheduledMeetings)
	route.GET("/users/:id/invites", handlers.GetUserInvites)

	route.GET("/events", handlers.GetEvents)
	route.GET("/events/:id", handlers.GetEventByID)
	route.POST("/events", handlers.CreateEvent)
	route.PATCH("/events/:id", handlers.UpdateEvent)

	route.POST("/events/:id/join", handlers.JoinEvent)

	route.GET("/events/:id/meetings", handlers.GetEventMeetings)
	route.GET("/events/:id/meetings/:meeting_id", handlers.GetEventMeetingByID)
	route.POST("/events/:id/meetings", handlers.CreateEventMeeting)

	route.GET("/meetings/:id/invites", handlers.GetMeetingInvites)
	route.GET("/meetings/:id/invites/:invite_id", handlers.GetMeetingInviteByID)
	route.POST("/meetings/:id/invites", handlers.CreateMeetingInvite)

	route.PATCH("/meetings/:id/schedule", handlers.ScheduleMeeting)

	route.PATCH("/invites/:id/rsvp", handlers.RespondToInvite)

	route.Run(":8085")
}
