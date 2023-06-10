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
	ID uint `gorm:"primaryKey"`

	Name     string
	Location string
	About    string

	Representatives []User
}

type User struct {
	ID uint `gorm:"primaryKey"`

	FirstName string
	LastName  string

	Location string
	About    string

	EMail    string
	Password string

	CompanyID uint

	OrganizedEvents   []Event   `gorm:"foreignKey:OrganizerID"`
	OrganizedMeetings []Meeting `gorm:"foreignKey:OrganizerID"`
	Invites           []Invite  `gorm:"foreignKey:InviteeID"`
}

type Event struct {
	ID uint `gorm:"primaryKey"`

	Name     string
	Location string
	Agenda   string

	StartDate time.Time
	EndDate   time.Time

	OrganizerID uint

	Participants []User `gorm:"many2many:event_participants;"`
	Meetings     []Meeting
}

type Meeting struct {
	ID uint `gorm:"primaryKey"`

	StartTime time.Time
	EndTime   time.Time

	Scheduled bool

	EventID     uint
	OrganizerID uint

	Invites []Invite
}

type Invite struct {
	ID uint `gorm:"primaryKey"`

	Status Status

	MeetingID uint
	InviteeID uint
}
