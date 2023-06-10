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

	Name     string `gorm:"not null" json:"name"`
	Location string `gorm:"not null" json:"location"`
	About    string `json:"about"`

	Representatives []User `json:"representatives"`
}

type User struct {
	ID uint `gorm:"primaryKey" json:"id"`

	FirstName string `gorm:"not null" json:"first_name"`
	LastName  string `gorm:"not null" json:"last_name"`

	Location string `json:"location"`
	About    string `json:"about"`

	EMail    string `gorm:"not null;unique" json:"e_mail"`
	Password string `gorm:"not null" json:"password"`

	CompanyID uint `json:"company_id"`

	OrganizedMeetings []Meeting `gorm:"foreignKey:OrganizerID" json:"organized_meetings"`
	Invites           []Invite  `gorm:"foreignKey:InviteeID" json:"-"`
}

type Event struct {
	ID uint `gorm:"primaryKey" json:"id"`

	Name     string `gorm:"not null" json:"name"`
	Location string `json:"location"`
	Agenda   string `json:"agenda"`

	StartDate time.Time `gorm:"not null" json:"start_date"`
	EndDate   time.Time `gorm:"not null" json:"end_date"`

	Participants []User    `gorm:"many2many:event_participants;" json:"participants"`
	Meetings     []Meeting `json:"-"`
}

type Meeting struct {
	ID uint `gorm:"primaryKey" json:"id"`

	StartTime time.Time `gorm:"not null" json:"start_time"`
	EndTime   time.Time `gorm:"not null" json:"end_time"`

	Scheduled bool `gorm:"not null;default:false" json:"scheduled"`

	EventID     uint `gorm:"not null" json:"event_id"`
	OrganizerID uint `gorm:"not null" json:"organizer_id"`

	Invites []Invite `json:"-"`
}

type Invite struct {
	ID uint `gorm:"primaryKey" json:"id"`

	Status Status `gorm:"not null;default:0" json:"status"`

	MeetingID uint `gorm:"not null" json:"meeting_id"`
	InviteeID uint `gorm:"not null" json:"invitee_id"`
}
