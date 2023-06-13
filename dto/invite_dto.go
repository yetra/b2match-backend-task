package dto

import "b2match/backend/models"

type NewInviteJSON struct {
	InviteeID uint `binding:"required" json:"invitee_id"`
}

type RSVPJSON struct {
	Status models.Status `binding:"required,min=1,max=2" json:"status"`
}
