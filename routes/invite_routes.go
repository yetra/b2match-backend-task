package routes

import (
	"b2match/backend/handlers"

	"github.com/gin-gonic/gin"
)

func addInviteGroup(router *gin.Engine) {
	inviteRouter := router.Group("/invites")
	{
		inviteRouter.GET("/:id", handlers.GetInviteByID)
		inviteRouter.PATCH("/:id/rsvp", handlers.RespondToInvite)
		inviteRouter.DELETE("/:id", handlers.DeleteInvite)
	}
}
