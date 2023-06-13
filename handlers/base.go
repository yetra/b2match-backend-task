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

// GET /<resource>
func getResources[R resource](c *gin.Context) {
	var resources []R
	database.DB.Find(&resources)

	c.JSON(http.StatusOK, resources)
}

// GET /<resource>/:id/<nested_resource>
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

// GET /<resource>/:id
func getResourceByID[R resource](c *gin.Context) {
	var resource R
	if err := findResourceByID(c, &resource, c.Param("id")); err != nil {
		return
	}

	c.JSON(http.StatusOK, resource)
}

// POST /<resource>
func createResource[R resource](c *gin.Context, resourceModel *R) {
	if err := database.DB.Create(resourceModel).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, *resourceModel)
}

// PATCH /<resource>/:id
func updateResource[R resource, J updateResourceJSON](c *gin.Context, resource *R, input *J) {
	if err := database.DB.Model(resource).Updates(input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": err.Error()})
		return
	}

	c.JSON(http.StatusOK, *resource)
}

// DELETE /<resource>/:id
func deleteResource[R resource](c *gin.Context, selectQuery interface{}) {
	var resource R
	if err := database.DB.First(&resource, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"errors": err.Error()})
		return
	}

	if selectQuery != nil {
		database.DB.Select(selectQuery).Delete(&resource)
	} else {
		database.DB.Delete(&resource)
	}

	c.Status(http.StatusNoContent)
}
