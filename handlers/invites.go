package handlers

import (
	"b2match/backend/dto"
	"b2match/backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetInviteByID godoc
//
// @Summary		 Get a single invite by id
// @Description	 Returns the invite whose ID value matches the id parameter.
// @Tags		 invites
// @Produce		 json
// @Param		 id path int true "Invite ID"
// @Success		 200	{object}	models.Invite
// @Failure		 404	{object}	dto.Error
// @Router		 /invites/{id} [get]
func GetInviteByID(c *gin.Context) {
	getResourceByID[models.Invite](c)
}

// RespondToInvite godoc
//
// @Summary		 Respond to an invite
// @Description	 Updates an invite's status with the request JSON. Returns the updated invite.
// @Tags		 invites
// @Produce		 json
// @Param		 id path int true "Invite ID"
// @Success		 200	{object}	models.Invite
// @Failure		 400	{object}	dto.Error
// @Failure		 404	{object}	dto.Error
// @Failure		 422	{object}	dto.Error
// @Router		 /invites/{id}/rsvp [patch]
func RespondToInvite(c *gin.Context) {
	var input dto.RSVPJSON
	if err := bindJSON(c, &input); err != nil {
		return
	}

	var invite models.Invite
	if err := findResourceByID(c, &invite, c.Param("id")); err != nil {
		return
	}

	if err := checkMeetingConflicts(invite.InviteeID, invite.MeetingID); err != nil {
		c.JSON(http.StatusUnprocessableEntity, dto.Error{Errors: err.Error()})
		return
	}

	updateResource(c, &invite, &input)
}

// DeleteInvite godoc
//
// @Summary      Delete an invite
// @Description  Deletes an invite specified by id.
// @Tags         invites
// @Accept       json
// @Produce      json
// @Param		 id	path int true "Invite ID"
// @Success      204  {object}  nil
// @Failure      404  {object}  dto.Error
// @Router       /invites/{id} [delete]
func DeleteInvite(c *gin.Context) {
	deleteResource[models.Invite](c, nil)
}
