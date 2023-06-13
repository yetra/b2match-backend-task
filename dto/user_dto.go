package dto

type NewUserJSON struct {
	FirstName string `binding:"required" json:"first_name"`
	LastName  string `binding:"required" json:"last_name"`

	Location string
	About    string

	EMail    string `binding:"required" json:"e_mail"`
	Password string `binding:"required"`

	CompanyID uint `binding:"required" json:"company_id"`
}

type UpdateUserJSON struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`

	Location string
	About    string

	Password string
}
