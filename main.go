package main

import (
	"b2match/backend/database"
	"b2match/backend/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	database.SetUpDB()

	route := gin.Default()

	route.GET("/companies", handlers.FindCompanies)
	route.GET("/companies/:id", handlers.FindCompanyByID)
	route.POST("/companies", handlers.CreateCompany)

	route.GET("/users", handlers.FindUsers)
	route.GET("/users/:id", handlers.FindUserByID)
	route.POST("/users", handlers.CreateUser)

	route.GET("/events", handlers.FindEvents)
	route.GET("/events/:id", handlers.FindEventByID)
	route.POST("/events", handlers.CreateEvent)

	route.POST("/events/:id/join", handlers.JoinEvent)

	route.GET("/events/:id/meetings", handlers.FindEventMeetings)
	route.GET("/events/:id/meetings/:meeting_id", handlers.FindEventMeetingByID)
	route.POST("/events/:id/meetings", handlers.CreateEventMeeting)

	route.GET("/meetings/:id/invites", handlers.FindMeetingInvites)
	route.GET("/meetings/:id/invites/:invite_id", handlers.FindMeetingInviteByID)
	route.POST("/meetings/:id/invites", handlers.CreateMeetingInvite)

	route.PATCH("/invites/:id/rsvp", handlers.RespondToInvite)

	route.Run(":8085")
}
