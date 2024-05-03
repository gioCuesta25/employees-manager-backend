package api

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gioCuesta25/employees-manager-backend/models"
	"github.com/gioCuesta25/employees-manager-backend/utils"
)

func (s *Server) createDepartment(ctx *gin.Context) {
	var body models.CreateDepartmentBody

	if err := ctx.ShouldBindJSON(&body); err != nil {
		utils.ErrorResponse(ctx, err, http.StatusBadRequest)
		return
	}

	query := `INSERT INTO departments
	(name, company_id)
	VALUES ($1, $2)
	RETURNING id, name, company_id, created_at, updated_at`

	row := s.db.QueryRow(query, body.Name, body.CompanyId)

	var department models.DepartmentsResponse

	err := row.Scan(&department.ID,
		&department.Name,
		&department.CompanyId,
		&department.CreatedAt,
		&department.UpdatedAt)

	if err != nil {
		utils.ErrorResponse(ctx, err, http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"department": department})
}

func (s *Server) getDepartmentsByCompany(ctx *gin.Context) {
	var params models.GetCompanyDepartmentsParams

	if err := ctx.ShouldBindUri(&params); err != nil {
		utils.ErrorResponse(ctx, err, http.StatusBadRequest)
		return
	}

	query := `SELECT id, name, company_id, created_at, updated_at FROM departments WHERE id = $1`

	rows, err := s.db.Query(query, params.CompanyId)

	departments := make([]*models.DepartmentsResponse, 0)

	for rows.Next() {
		d, err := scanRowsIntoDepartment(rows)
		if err != nil {
			utils.ErrorResponse(ctx, err, http.StatusInternalServerError)
			return
		}
		departments = append(departments, d)
	}

	ctx.JSON(http.StatusOK, gin.H{"items": departments})

	if err != nil {
		utils.ErrorResponse(ctx, err, http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"items": departments})
}

func (s *Server) updateDepartment(ctx *gin.Context) {
	var body models.CreateDepartmentBody
	var params models.GetDepartmentsParams

	if err := ctx.ShouldBindJSON(&body); err != nil {
		utils.ErrorResponse(ctx, err, http.StatusBadRequest)
		return
	}

	if err := ctx.ShouldBindUri(&params); err != nil {
		utils.ErrorResponse(ctx, err, http.StatusBadRequest)
		return
	}

	query := `UPDATE departments
	SET name = $1, updated_at = $2
	WHERE id = $3
	RETURNING id, name, company_id, created_at, updated_at`

	row := s.db.QueryRow(query, body.Name, body.CompanyId, time.Now(), params.ID)

	var department models.DepartmentsResponse

	err := row.Scan(&department.ID,
		&department.Name,
		&department.CompanyId,
		&department.CreatedAt,
		&department.UpdatedAt)

	if err != nil {
		utils.ErrorResponse(ctx, err, http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"department": department})
}

func (s *Server) deleteDepartment(ctx *gin.Context) {
	var params models.GetDepartmentsParams

	if err := ctx.ShouldBindUri(&params); err != nil {
		utils.ErrorResponse(ctx, err, http.StatusBadRequest)
		return
	}

	query := `DELETE FROM departments WHERE id = $1`

	_, err := s.db.Query(query, params.ID)

	if err != nil {
		utils.ErrorResponse(ctx, err, http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Department deleted successfully"})

}

func scanRowsIntoDepartment(rows *sql.Rows) (*models.DepartmentsResponse, error) {
	department := new(models.DepartmentsResponse)

	err := rows.Scan(
		&department.ID,
		&department.Name,
		&department.CompanyId,
		&department.CreatedAt,
		&department.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return department, nil
}
