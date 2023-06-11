package models

type User struct {
	ID uint `gorm:"primaryKey;<-:create" json:"id"`

	FirstName string `gorm:"not null" json:"first_name" binding:"required"`
	LastName  string `gorm:"not null" json:"last_name" binding:"required"`

	Location string `json:"location"`
	About    string `json:"about"`

	EMail    string `gorm:"not null;unique" json:"e_mail" binding:"required"`
	Password string `gorm:"not null" json:"password" binding:"required"`

	CompanyID uint `gorm:"not null;<-:create" json:"company_id" binding:"required"`

	Events []Event `gorm:"many2many:event_participants;" json:"events"`

	OrganizedMeetings []Meeting `gorm:"foreignKey:OrganizerID" json:"organized_meetings"`
	Invites           []Invite  `gorm:"foreignKey:InviteeID" json:"-"`
}
