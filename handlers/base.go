package handlers

import (
	"b2match/backend/database"
	"b2match/backend/models"
	"net/http"

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

func findResourceByID[R resource](c *gin.Context, resource *R, id interface{}) error {
	if err := database.DB.First(resource, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"errors": err.Error()})
		return err
	}

	return nil
}

func findNestedResources[R, RNested resource](c *gin.Context, resource *R, nestedResources *[]RNested, assocName string) error {
	err := database.DB.Model(resource).Association(assocName).Find(nestedResources)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"errors": err.Error()})
		return err
	}

	return nil
}

func bindJSON[J inputJSON](c *gin.Context, inputJSON *J) error {
	if err := c.ShouldBindJSON(inputJSON); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return err
	}

	return nil
}

func getResources[R resource](c *gin.Context) {
	var resources []R
	database.DB.Find(&resources)

	c.JSON(http.StatusOK, resources)
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

	c.JSON(http.StatusOK, nestedResources)
}

func getResourceByID[R resource](c *gin.Context) {
	var resource R
	if err := findResourceByID(c, &resource, c.Param("id")); err != nil {
		return
	}

	c.JSON(http.StatusOK, resource)
}

func createResource[R resource](c *gin.Context, resourceModel *R) {
	if err := database.DB.Create(resourceModel).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, *resourceModel)
}

func updateResource[R resource, J updateResourceJSON](c *gin.Context, resource *R, input *J) {
	if err := database.DB.Model(resource).Updates(input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": err.Error()})
		return
	}

	c.JSON(http.StatusOK, *resource)
}

func deleteResource[R resource](c *gin.Context, selectQuery interface{}) {
	var resource R
	if err := database.DB.First(&resource, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"errors": err.Error()})
		return
	}

	database.DB.Select(selectQuery).Delete(&resource)

	c.Status(http.StatusNoContent)
}
