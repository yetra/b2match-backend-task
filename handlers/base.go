package handlers

import (
	"b2match/backend/database"
	"b2match/backend/models"
	"net/http"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
)

type resource interface {
	models.Company | models.User | models.Event | models.Meeting | models.Invite
}

func FindResources[r resource](c *gin.Context) {
	var resources []r
	database.DB.Find(&resources)

	resourcesName := strings.ToLower(reflect.TypeOf(resources).Elem().Name())

	c.JSON(http.StatusOK, gin.H{resourcesName: resources})
}
