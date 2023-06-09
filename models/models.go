package models

import (
	"time"

	"gorm.io/gorm"
)

type Company struct {
	gorm.Model

	Name     string
	Location string
	About    string
}

type User struct {
	gorm.Model

	FirstName string
	LastName  string

	Location string
	About    string

	EMail    string
	Password string
}

type Event struct {
	gorm.Model

	Name     string
	Location string
	Agenda   string

	StartDate time.Time
	EndDate   time.Time
}

type Meeting struct {
	gorm.Model

	StartTime time.Time
	EndTime   time.Time

	Scheduled bool
}

type Invite struct {
	gorm.Model

	Status string
}
