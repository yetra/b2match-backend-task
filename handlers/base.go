package handlers

import (
	"b2match/backend/database"
	"b2match/backend/dto"
	"net/http"

	"github.com/gin-gonic/gin"
)

type newResourceJSON interface {
	dto.NewCompanyJSON | dto.NewUserJSON | dto.NewEventJSON | dto.NewMeetingJSON | dto.NewInviteJSON
}
type updateResourceJSON interface {
	dto.UpdateCompanyJSON | dto.UpdateUserJSON | dto.UpdateEventJSON | dto.RSVPJSON
}
type inputJSON interface {
	newResourceJSON | updateResourceJSON | dto.JoinEventJSON
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
