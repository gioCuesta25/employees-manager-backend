package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gioCuesta25/employees-manager-backend/models"
	"github.com/gioCuesta25/employees-manager-backend/utils"
)

func (s *Server) createCompany(ctx *gin.Context) {
	var body models.CreateCompanyBody

	if err := ctx.ShouldBindJSON(&body); err != nil {
		utils.ErrorResponse(ctx, err, http.StatusBadRequest)
		return
	}

	query := `INSERT INTO companies
			(name, owner, address, phone, email) 
			VALUES ($1, $2, $3, $4, $5)
			RETURNING id, name, owner, address, phone, email, created_at, updated_at`

	row := s.db.QueryRow(query, body.Name, body.Owner, body.Address, body.Phone, body.Email)
	var newCompany models.CompanyResponse
	err := row.Scan(&newCompany.ID,
		&newCompany.Name,
		&newCompany.Owner,
		&newCompany.Address,
		&newCompany.Phone,
		&newCompany.Email,
		&newCompany.CreatedAt,
		&newCompany.UpdatedAt)

	if err != nil {
		utils.ErrorResponse(ctx, err, http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"company": newCompany})
}

func (s *Server) getCompany(ctx *gin.Context) {
	var params models.GetCompanyParams

	if err := ctx.ShouldBindUri(&params); err != nil {
		utils.ErrorResponse(ctx, err, http.StatusBadRequest)
		return
	}

	query := "SELECT id, name, owner, address, phone, email, created_at, updated_at FROM companies WHERE id = $1"
	row := s.db.QueryRow(query, params.ID)

	var company models.CompanyResponse
	err := row.Scan(&company.ID,
		&company.Name,
		&company.Owner,
		&company.Address,
		&company.Phone,
		&company.Email,
		&company.CreatedAt,
		&company.UpdatedAt)

	if err != nil {
		utils.ErrorResponse(ctx, err, http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"company": company})
}

func (s *Server) deleteCompany(ctx *gin.Context) {
	var params models.GetCompanyParams

	if err := ctx.ShouldBindUri(&params); err != nil {
		utils.ErrorResponse(ctx, err, http.StatusBadRequest)
		return
	}

	query := "DELETE FROM companies WHERE id = $1"
	_, err := s.db.Query(query, params.ID)

	if err != nil {
		utils.ErrorResponse(ctx, err, http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Company successfully deleted"})
}

func (s *Server) updateCompany(ctx *gin.Context) {
	var params models.GetCompanyParams
	var body models.CreateCompanyBody

	if err := ctx.ShouldBindUri(&params); err != nil {
		utils.ErrorResponse(ctx, err, http.StatusBadRequest)
		return
	}

	if err := ctx.ShouldBindJSON(&body); err != nil {
		utils.ErrorResponse(ctx, err, http.StatusBadRequest)
		return
	}

	query := `UPDATE companies
			name = $1,
			owner = $2,
			address = $3,
			phone = $4,
			email = $5
			WHERE id = $6
			RETURNING id, name, owner, address, phone, email, created_at, updated_at`

	row := s.db.QueryRow(query, body.Name, body.Owner, body.Address, body.Phone, body.Email, params.ID)
	var newCompany models.CompanyResponse
	err := row.Scan(&newCompany.ID,
		&newCompany.Name,
		&newCompany.Owner,
		&newCompany.Address,
		&newCompany.Phone,
		&newCompany.Email,
		&newCompany.CreatedAt,
		&newCompany.UpdatedAt)

	if err != nil {
		utils.ErrorResponse(ctx, err, http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"company": newCompany})

}
