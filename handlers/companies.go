package handlers

import (
	"b2match/backend/dto"
	"b2match/backend/models"

	"github.com/gin-gonic/gin"
)

// GetCompanies godoc
//
// @Summary		 Get companies
// @Description	 Responds with a list of all companies as JSON.
// @Tags		 companies
// @Produce		 json
// @Success		 200	{array}		models.Company
// @Router		 /companies [get]
func GetCompanies(c *gin.Context) {
	getResources[models.Company](c)
}

// GetCompanyByID godoc
//
// @Summary		 Get a single company by id
// @Description	 Returns the company whose ID value matches the id parameter.
// @Tags		 companies
// @Produce		 json
// @Param		 id path int true "Company ID"
// @Success		 200	{object}	models.Company
// @Failure		 404	{object}	gin.H
// @Router		 /companies/{id} [get]
func GetCompanyByID(c *gin.Context) {
	getResourceByID[models.Company](c)
}

// CreateCompany godoc
//
// @Summary      Create a new company
// @Description  Creates a company from the input JSON and stores it in the database. Returns the new company.
// @Tags         companies
// @Accept       json
// @Produce      json
// @Success      201 	{object}	models.Company
// @Failure      400 	{object}	gin.H
// @Failure      500 	{object}	gin.H
// @Router       /companies [post]
func CreateCompany(c *gin.Context) {
	var input dto.NewCompanyJSON
	if err := bindJSON(c, &input); err != nil {
		return
	}

	company := models.Company{
		Name:     input.Name,
		Location: input.Location,
		About:    input.About,
	}

	createResource(c, &company)
}

// UpdateCompany godoc
//
// @Summary      Update an existing company
// @Description  Updates a company with the input JSON. Returns the updated company.
// @Tags         companies
// @Accept       json
// @Produce      json
// @Param		 id	path int true "Company ID"
// @Success      200	{object}	models.Company
// @Failure      400	{object}	gin.H
// @Failure      404	{object}	gin.H
// @Failure      500	{object}	gin.H
// @Router       /companies/{id} [patch]
func UpdateCompany(c *gin.Context) {
	var company models.Company
	if err := findResourceByID(c, &company, c.Param("id")); err != nil {
		return
	}

	var input dto.UpdateCompanyJSON
	if err := bindJSON(c, &input); err != nil {
		return
	}

	updateResource(c, &company, &input)
}

// DeleteCompany godoc
//
// @Summary      Delete a company
// @Description  Deletes a company and its representatives.
// @Tags         companies
// @Accept       json
// @Produce      json
// @Param		 id	path int true "Company ID"
// @Success      204  {object}  nil
// @Failure      404  {object}  gin.H
// @Router       /companies/{id} [delete]
func DeleteCompany(c *gin.Context) {
	deleteResource[models.Company](c, "Representatives")
}

// GetCompanyRepresentatives godoc
//
// @Summary      Get company representatives
// @Description  Responds with a list of company representatives as JSON.
// @Tags         companies
// @Produce      json
// @Param		 id	path int true "Company ID"
// @Success      200	{array}		models.User
// @Failure      404	{object}	gin.H
// @Failure      500	{object}	gin.H
// @Router       /companies/{id}/representatives [get]
func GetCompanyRepresentatives(c *gin.Context) {
	getNestedResources[models.Company, models.User](c, "Representatives")
}
