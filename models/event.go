package models

import "time"

type Event struct {
	ID uint `gorm:"primaryKey;<-:create" json:"id"`

	Name     string `gorm:"not null" json:"name" binding:"required"`
	Location string `json:"location"`
	Agenda   string `json:"agenda"`

	StartDate time.Time `gorm:"not null" json:"start_date" binding:"required,ltefield=EndDate"`
	EndDate   time.Time `gorm:"not null" json:"end_date" binding:"required"`

	Participants []User    `gorm:"many2many:event_participants;" json:"-"`
	Meetings     []Meeting `gorm:"constraint:OnDelete:CASCADE" json:"-"`
}
