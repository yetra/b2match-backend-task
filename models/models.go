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

	OrganizedEvents   []Event
	OrganizedMeetings []Meeting
}

type Event struct {
	gorm.Model

	Name     string
	Location string
	Agenda   string

	StartDate time.Time
	EndDate   time.Time

	OrganizerID int
}

type Meeting struct {
	gorm.Model

	StartTime time.Time
	EndTime   time.Time

	Scheduled bool

	OrganizerID int
}

type Invite struct {
	gorm.Model

	Status Status
}
