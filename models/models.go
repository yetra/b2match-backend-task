package models

import (
	"time"

	"gorm.io/gorm"
)

type Status int

const (
	Pending  Status = iota
	Accepted Status = iota
	Declined Status = iota
)

type Company struct {
	gorm.Model

	Name     string
	Location string
	About    string

	Users []User
}

type User struct {
	gorm.Model

	FirstName string
	LastName  string

	Location string
	About    string

	EMail    string
	Password string

	CompanyID int

	OrganizedEvents   []Event   `gorm:"foreignKey:OrganizerID"`
	OrganizedMeetings []Meeting `gorm:"foreignKey:OrganizerID"`
	Invites           []Invite  `gorm:"foreignKey:InviteeID"`
}

type Event struct {
	gorm.Model

	Name     string
	Location string
	Agenda   string

	StartDate time.Time
	EndDate   time.Time

	OrganizerID int

	Participants []User `gorm:"many2many:event_participants;"`
	Meetings     []Meeting
}

type Meeting struct {
	gorm.Model

	StartTime time.Time
	EndTime   time.Time

	Scheduled bool

	EventID     int
	OrganizerID int

	Invites []Invite
}

type Invite struct {
	gorm.Model

	Status Status

	MeetingID int
	InviteeID uint
}
