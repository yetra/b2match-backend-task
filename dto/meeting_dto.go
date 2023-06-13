package dto

import "time"

type NewMeetingJSON struct {
	StartTime time.Time `binding:"required,ltefield=EndTime" json:"start_time"`
	EndTime   time.Time `binding:"required" json:"end_time"`

	OrganizerID uint `binding:"required" json:"organizer_id"`
}
