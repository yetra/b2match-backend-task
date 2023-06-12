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

func findResourceByID[R resource](c *gin.Context, resource *R, id interface{}) error {
	if err := database.DB.First(resource, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return err
	}

	return nil
}

func getResources[R resource](c *gin.Context) {
	var resources []R
	database.DB.Find(&resources)

	resourcesName := strings.ToLower(getTypeName(resources))

	c.JSON(http.StatusOK, gin.H{resourcesName: resources})
}

func getNestedResources[R, RNested resource](c *gin.Context, assocName string) {
	var resource R
	var nestedResources []RNested

	if err := findResourceByID(c, &resource, c.Param("id")); err != nil {
		return
	}

	err := database.DB.Model(&resource).Association(assocName).Find(&nestedResources)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	nestedResourcesName := strings.ToLower(getTypeName(nestedResources))

	c.JSON(http.StatusOK, gin.H{nestedResourcesName: nestedResources})
}

func getResourceByID[R resource](c *gin.Context) {
	var resource R

	if err := findResourceByID(c, &resource, c.Param("id")); err != nil {
		return
	}

	resourceName := strings.ToLower(getTypeName(resource))

	c.JSON(http.StatusOK, gin.H{resourceName: resource})
}
