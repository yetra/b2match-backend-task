package handlers

import (
	"b2match/backend/database"
	"b2match/backend/dto"
	"b2match/backend/models"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GET /meetings/:id/invites
func GetMeetingInvites(c *gin.Context) {
	getNestedResources[models.Meeting, models.Invite](c, "Invites")
}

// POST /meetings/:id/invites
func CreateMeetingInvite(c *gin.Context) {
	var input dto.NewInviteJSON
	if err := bindJSON(c, &input); err != nil {
		return
	}

	var meeting models.Meeting
	if err := findResourceByID(c, &meeting, c.Param("id")); err != nil {
		return
	}

	var invitee models.User
	if err := findResourceByID(c, &invitee, input.InviteeID); err != nil {
		return
	}

	if err := checkInviteeIsNotOrganizer(invitee.ID, meeting.OrganizerID); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"errors": err.Error()})
		return
	}
	if err := checkInviteeIsAParticipant(invitee.ID, meeting.EventID); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"errors": err.Error()})
		return
	}
	if err := checkInviteeNotAlreadyInvited(invitee.ID, meeting.ID); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"errors": err.Error()})
		return
	}

	invite := models.Invite{
		Status: models.Pending,

		MeetingID: meeting.ID,
		InviteeID: invitee.ID,
	}

	createResource(c, &invite)
}

// GET /invites/:id
func GetInviteByID(c *gin.Context) {
	getResourceByID[models.Meeting](c)
}

// PATCH /invites/:id/rsvp
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
		c.JSON(http.StatusUnprocessableEntity, gin.H{"errors": err.Error()})
		return
	}

	updateResource(c, &invite, &input)
}

// DELETE /invites/:id
func DeleteInvite(c *gin.Context) {
	deleteResource[models.Invite](c, nil)
}

func checkInviteeIsNotOrganizer(inviteeID uint, organizerID uint) error {
	if inviteeID == organizerID {
		return errors.New("the invitee is the meeting organizer")
	}

	return nil
}

func checkInviteeIsAParticipant(inviteeID uint, eventID uint) error {
	var event models.Event

	err := database.DB.Preload("Participants").First(&event, eventID).Error
	if err != nil {
		return err
	}

	for _, participant := range event.Participants {
		if participant.ID == inviteeID {
			return nil
		}
	}

	return errors.New("the invitee is not an event participant")
}

func checkInviteeNotAlreadyInvited(inviteeID uint, meetingID uint) error {
	var invites []models.Invite

	err := database.DB.Where("meeting_id = ?", meetingID).Find(&invites).Error
	if err != nil {
		return err
	}

	for _, invite := range invites {
		if inviteeID == invite.InviteeID {
			return errors.New("invitee already invited to meeting")
		}
	}

	return nil
}

func checkMeetingConflicts(inviteeID uint, meetingID uint) error {
	var acceptedInvites []models.Invite

	whereClause := "invitee_id = ? AND status = ?"
	err := database.DB.Find(&acceptedInvites, whereClause, inviteeID, models.Accepted).Error
	if err != nil {
		return err
	}

	var meeting models.Meeting
	if err := database.DB.First(&meeting, meetingID).Error; err != nil {
		return err
	}

	for _, invite := range acceptedInvites {
		var acceptedMeeting models.Meeting
		if err := database.DB.First(&acceptedMeeting, invite.MeetingID).Error; err != nil {
			return err
		}

		if meeting.StartTime.Before(acceptedMeeting.EndTime) &&
			meeting.EndTime.After(acceptedMeeting.StartTime) {
			return errors.New("this meeting is in conflict with already accepted meetings")
		}
	}

	return nil
}
