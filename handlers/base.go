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

type newResourceJSON interface {
	newCompanyJSON | newUserJSON | newEventJSON | newMeetingJSON | newInviteJSON
}

type updateResourceJSON interface {
	updateCompanyJSON | updateUserJSON | updateEventJSON | rsvpJSON
}

type inputJSON interface {
	newResourceJSON | updateResourceJSON | joinEventJSON
}

func typeName(variable interface{}) string {
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

func findNestedResources[R, RNested resource](c *gin.Context, resource *R, nestedResources *[]RNested, assocName string) error {
	err := database.DB.Model(&resource).Association(assocName).Find(&nestedResources)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return err
	}

	return nil
}

func bindJSON[J inputJSON](c *gin.Context, inputJSON *J) error {
	if err := c.ShouldBindJSON(&inputJSON); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return err
	}

	return nil
}

func getResources[R resource](c *gin.Context) {
	var resources []R
	database.DB.Find(&resources)

	resourcesName := strings.ToLower(typeName(resources))

	c.JSON(http.StatusOK, gin.H{resourcesName: resources})
}

func getNestedResources[R, RNested resource](c *gin.Context, assocName string) {
	var resource R
	if err := findResourceByID(c, &resource, c.Param("id")); err != nil {
		return
	}

	var nestedResources []RNested
	if err := findNestedResources(c, &resource, &nestedResources, assocName); err != nil {
		return
	}

	nestedResourcesName := strings.ToLower(typeName(nestedResources))

	c.JSON(http.StatusOK, gin.H{nestedResourcesName: nestedResources})
}

func getResourceByID[R resource](c *gin.Context) {
	var resource R
	if err := findResourceByID(c, &resource, c.Param("id")); err != nil {
		return
	}

	resourceName := strings.ToLower(typeName(resource))

	c.JSON(http.StatusOK, gin.H{resourceName: resource})
}

func createResource[R resource](c *gin.Context, resourceModel *R) {
	if err := database.DB.Create(&resourceModel).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resourceName := strings.ToLower(typeName(resourceModel))

	c.JSON(http.StatusCreated, gin.H{resourceName: resourceModel})
}

func updateResource[R resource, J updateResourceJSON](c *gin.Context, resource *R, input *J) {
	if err := database.DB.Model(&resource).Updates(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resourceName := strings.ToLower(typeName(resource))

	c.JSON(http.StatusOK, gin.H{resourceName: resource})
}
