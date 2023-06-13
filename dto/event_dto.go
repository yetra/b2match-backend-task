package dto

import "time"

type NewEventJSON struct {
	Name     string `binding:"required"`
	Location string
	Agenda   string

	StartDate time.Time `binding:"required,ltefield=EndDate" json:"start_date"`
	EndDate   time.Time `binding:"required" json:"end_date"`
}

type UpdateEventJSON struct {
	Agenda string
}

type JoinEventJSON struct {
	ID uint `binding:"required" json:"id"`
}
