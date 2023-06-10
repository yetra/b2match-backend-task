package models

import (
	"time"
)

type Status int

const (
	Pending  Status = iota
	Accepted Status = iota
	Declined Status = iota
)

type Company struct {
	ID uint `gorm:"primaryKey" json:"id"`

	Name     string `gorm:"not null" json:"name" binding:"required"`
	Location string `gorm:"not null" json:"location"`
	About    string `json:"about"`

	Representatives []User `json:"representatives"`
}

type User struct {
	ID uint `gorm:"primaryKey" json:"id"`

	FirstName string `gorm:"not null" json:"first_name" binding:"required"`
	LastName  string `gorm:"not null" json:"last_name" binding:"required"`

	Location string `json:"location"`
	About    string `json:"about"`

	EMail    string `gorm:"not null;unique" json:"e_mail" binding:"required"`
	Password string `gorm:"not null" json:"password" binding:"required"`

	CompanyID uint `json:"company_id" binding:"required"`

	OrganizedMeetings []Meeting `gorm:"foreignKey:OrganizerID" json:"organized_meetings"`
	Invites           []Invite  `gorm:"foreignKey:InviteeID" json:"-"`
}

type Event struct {
	ID uint `gorm:"primaryKey" json:"id"`

	Name     string `gorm:"not null" json:"name" binding:"required"`
	Location string `json:"location"`
	Agenda   string `json:"agenda"`

	StartDate time.Time `gorm:"not null" json:"start_date" binding:"required,ltefield=EndDate"`
	EndDate   time.Time `gorm:"not null" json:"end_date" binding:"required"`

	Participants []User    `gorm:"many2many:event_participants;" json:"participants"`
	Meetings     []Meeting `json:"-"`
}

type Meeting struct {
	ID uint `gorm:"primaryKey" json:"id"`

	StartTime time.Time `gorm:"not null" json:"start_time" binding:"required,ltefield=EndTime"`
	EndTime   time.Time `gorm:"not null" json:"end_time" binding:"required"`

	Scheduled bool `gorm:"not null;default:false" json:"scheduled"`

	EventID     uint `gorm:"not null" json:"event_id"`
	OrganizerID uint `gorm:"not null" json:"organizer_id" binding:"required"`

	Invites []Invite `json:"-"`
}

type Invite struct {
	ID uint `gorm:"primaryKey" json:"id"`

	Status Status `gorm:"not null;default:0" json:"status"`

	MeetingID uint `gorm:"not null" json:"meeting_id"`
	InviteeID uint `gorm:"not null" json:"invitee_id" binding:"required"`
}
