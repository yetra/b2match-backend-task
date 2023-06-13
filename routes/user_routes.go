package routes

import (
	"b2match/backend/handlers"

	"github.com/gin-gonic/gin"
)

func addUserGroup(router *gin.Engine) {
	userRouter := router.Group("/users")
	{
		userRouter.GET("", handlers.GetUsers)
		userRouter.GET("/:id", handlers.GetUserByID)
		userRouter.POST("", handlers.CreateUser)
		userRouter.PATCH("/:id", handlers.UpdateUser)
		userRouter.DELETE("/:id", handlers.DeleteUser)

		userRouter.GET("/:id/scheduled-meetings", handlers.GetUserScheduledMeetings)
		userRouter.GET("/:id/invites", handlers.GetUserInvites)
	}
}
