package handlers

import (
	"b2match/backend/database"
	"b2match/backend/dto"
	"b2match/backend/models"
	"errors"
)

func checkNewMeetingIsDuringEvent(meeting dto.NewMeetingJSON, event models.Event) error {
	if meeting.StartTime.After(event.StartDate) && meeting.EndTime.Before(event.EndDate) {
		return nil
	}

	return errors.New("meeting must happen during the event")
}

func checkMeetingNotAlreadyScheduled(meeting models.Meeting) error {
	if meeting.Scheduled {
		return errors.New("meeting is already scheduled")
	}

	return nil
}

func checkInviteeIsNotOrganizer(inviteeID uint, organizerID uint) error {
	if inviteeID == organizerID {
		return errors.New("cannot invite the meeting organizer")
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

	return errors.New("invitee is not an event participant")
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
	err := database.DB.Find(&acceptedInvites, whereClause, inviteeID, models.StatusAccepted).Error
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
