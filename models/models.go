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

	Name     string `json:"name"`
	Location string `json:"location"`
	About    string `json:"about"`

	Representatives []User `json:"representatives"`
}

type User struct {
	ID uint `gorm:"primaryKey" json:"id"`

	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`

	Location string `json:"location"`
	About    string `json:"about"`

	EMail    string `json:"e_mail"`
	Password string `json:"password"`

	CompanyID uint `json:"company_id"`

	OrganizedEvents   []Event   `gorm:"foreignKey:OrganizerID" json:"organized_events"`
	OrganizedMeetings []Meeting `gorm:"foreignKey:OrganizerID" json:"organized_meetings"`
	Invites           []Invite  `gorm:"foreignKey:InviteeID" json:"invites"`
}

type Event struct {
	ID uint `gorm:"primaryKey" json:"id"`

	Name     string `json:"name"`
	Location string `json:"location"`
	Agenda   string `json:"agenda"`

	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`

	OrganizerID uint `json:"organizer_id"`

	Participants []User    `gorm:"many2many:event_participants;" json:"participants"`
	Meetings     []Meeting `json:"meetings"`
}

type Meeting struct {
	ID uint `gorm:"primaryKey" json:"id"`

	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`

	Scheduled bool `json:"scheduled"`

	EventID     uint `json:"event_id"`
	OrganizerID uint `json:"organizer_id"`

	Invites []Invite `json:"invites"`
}

type Invite struct {
	ID uint `gorm:"primaryKey" json:"id"`

	Status Status `json:"status"`

	MeetingID uint `json:"meeting_id"`
	InviteeID uint `json:"invitee_id"`
}
