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

func getTypeName(variable interface{}) string {
	t := reflect.TypeOf(variable)

	if t.Kind() == reflect.Slice || t.Kind() == reflect.Ptr {
		return t.Elem().Name()
	} else {
		return t.Name()
	}
}

func findResources[R resource](c *gin.Context) {
	var resources []R
	database.DB.Find(&resources)

	resourcesName := strings.ToLower(getTypeName(resources))

	c.JSON(http.StatusOK, gin.H{resourcesName: resources})
}

func findResourceByID[r resource](c *gin.Context) {
	var resource r

	id := c.Param("id")

	if err := database.DB.First(&resource, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	resourceName := strings.ToLower(getTypeName(resource))

	c.JSON(http.StatusOK, gin.H{resourceName: resource})
}
