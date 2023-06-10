package models

import "time"

type Meeting struct {
	ID uint `gorm:"primaryKey;<-:create" json:"id"`

	StartTime time.Time `gorm:"not null" json:"start_time" binding:"required,ltefield=EndTime"`
	EndTime   time.Time `gorm:"not null" json:"end_time" binding:"required"`

	Scheduled bool `gorm:"not null;default:false;<-:update" json:"scheduled"`

	EventID     uint `gorm:"not null;<-:create" json:"event_id"`
	OrganizerID uint `gorm:"not null;<-:create" json:"organizer_id" binding:"required"`

	Invites []Invite `gorm:"->" json:"-"`
}
