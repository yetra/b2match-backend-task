package models

type Company struct {
	ID uint `gorm:"primaryKey;<-:create" json:"id"`

	Name     string `gorm:"not null" json:"name" binding:"required"`
	Location string `gorm:"not null" json:"location"`
	About    string `json:"about"`

	Representatives []User `json:"representatives"`
}
