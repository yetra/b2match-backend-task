package models

type Status int

const (
	StatusPending Status = iota
	StatusAccepted
	StatusDeclined
)

type Invite struct {
	ID uint `gorm:"primaryKey;<-:create" json:"id"`

	Status Status `gorm:"not null;default:0;<-:update" json:"status"`

	MeetingID uint `gorm:"not null;<-:create" json:"meeting_id"`
	InviteeID uint `gorm:"not null;<-:create" json:"invitee_id" binding:"required"`
}
