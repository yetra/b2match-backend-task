package routes

import "github.com/gin-gonic/gin"

// CreateRouter creates a new gin.Engine instance and adds routes to it.
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
