package routes

import (
	"b2match/backend/handlers"

	"github.com/gin-gonic/gin"
)

func addMeetingGroup(router *gin.Engine) {
	meetingRouter := router.Group("/meetings")
	{
		meetingRouter.GET("/:id", handlers.GetMeetingByID)
		meetingRouter.PATCH("/:id/schedule", handlers.ScheduleMeeting)
		meetingRouter.DELETE("/:id", handlers.DeleteMeeting)

		meetingRouter.GET("/:id/invites", handlers.GetMeetingInvites)
		meetingRouter.POST("/:id/invites", handlers.CreateMeetingInvite)
	}
}
