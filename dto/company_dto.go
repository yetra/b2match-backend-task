package dto

type NewCompanyJSON struct {
	Name     string `binding:"required"`
	Location string `binding:"required"`
	About    string
}

type UpdateCompanyJSON struct {
	Location string
	About    string
}
