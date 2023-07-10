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

func findResourceByID[R resource](c *gin.Context, id interface{}) (resource R, err error) {
	if err := database.DB.First(&resource, id).Error; err != nil {
		c.JSON(http.StatusNotFound, dto.Error{Errors: err.Error()})
		return resource, err
	}

	return resource, nil
}

func findNestedResources[RNested, R resource](c *gin.Context, resource *R, assocName string) (nestedResources []RNested, err error) {
	err = database.DB.Model(resource).Association(assocName).Find(&nestedResources)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.Error{Errors: err.Error()})
		return nestedResources, err
	}

	return nestedResources, nil
}

func bindJSON[J inputJSON](c *gin.Context) (input J, err error) {
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error{Errors: err.Error()})
		return input, err
	}

	return input, nil
}
