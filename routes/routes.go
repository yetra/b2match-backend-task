package routes

import "github.com/gin-gonic/gin"

func CreateRouter() *gin.Engine {
	router := gin.Default()
	addRoutes(router)

	return router
}

func addRoutes(router *gin.Engine) {
	addCompanyGroup(router)
	addUserGroup(router)
	addEventGroup(router)
	addMeetingGroup(router)
	addInviteGroup(router)
}
