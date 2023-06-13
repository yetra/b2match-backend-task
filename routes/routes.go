package routes

import "github.com/gin-gonic/gin"

func AddRoutes(router *gin.Engine) {
	addCompanyGroup(router)
	addUserGroup(router)
	addEventGroup(router)
	addMeetingGroup(router)
	addInviteGroup(router)
}
