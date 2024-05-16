package api

import (
	"fmt"
	"math"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gioCuesta25/employees-manager-backend/models"
	"github.com/gioCuesta25/employees-manager-backend/utils"
)

func (s *Server) createEmployee(ctx *gin.Context) {
	var body models.CreateEmployeeBody

	if err := ctx.ShouldBindJSON(&body); err != nil {
		utils.ErrorResponse(ctx, err, http.StatusBadRequest)
		return
	}

	query := `INSERT INTO employees (name,
		last_name,
		phone_number,
		email,
		id_type,
		id_number,
		admission_date,
		salary,
		position_id,
		department_id,
		company_id,
		picture_url)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id, name, last_name, phone_number, email, id_type, id_number, admission_date, salary, position_id, department_id,company_id, picture_url,created_at, updated_at`

	row := s.db.QueryRow(
		query,
		body.Name,
		body.LastName,
		body.PhoneNumber,
		body.Email,
		body.IdType,
		body.IdNumber,
		body.AdmissionDate,
		body.Salary,
		body.PositionId,
		body.DepartmentId,
		body.CompanyId,
		body.PictureUrl)

	var employee models.EmployeeResponse

	err := row.Scan(
		&employee.ID,
		&employee.Name,
		&employee.LastName,
		&employee.PhoneNumber,
		&employee.Email,
		&employee.IdType,
		&employee.IdNumber,
		&employee.AdmissionDate,
		&employee.Salary,
		&employee.PositionId,
		&employee.DepartmentId,
		&employee.CompanyId,
		&employee.PictureUrl,
		&employee.CreatedAt,
		&employee.UpdatedAt)

	if err != nil {
		utils.ErrorResponse(ctx, err, http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"employee": employee})
}

func (s *Server) getEmployeeById(ctx *gin.Context) {
	var params models.GetEmployeeParams

	if err := ctx.ShouldBindUri(&params); err != nil {
		utils.ErrorResponse(ctx, err, http.StatusBadRequest)
		return
	}

	query := `
	SELECT
		id,
		name,
		last_name,
		phone_number,
		email,
		id_type,
		id_number,
		admission_date,
		salary,
		position_id,
		department_id,
		company_id,
		picture_url,
		created_at,
		updated_at
	FROM
		employees
	WHERE
		id = $1
	`

	row := s.db.QueryRow(query, params.ID)

	var employee models.EmployeeResponse

	err := row.Scan(
		&employee.ID,
		&employee.Name,
		&employee.LastName,
		&employee.PhoneNumber,
		&employee.Email,
		&employee.IdType,
		&employee.IdNumber,
		&employee.AdmissionDate,
		&employee.Salary,
		&employee.PositionId,
		&employee.DepartmentId,
		&employee.CompanyId,
		&employee.CreatedAt,
		&employee.UpdatedAt)

	if err != nil {
		utils.ErrorResponse(ctx, err, http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"employee": employee})
}

func (s *Server) listCompanyEmployees(ctx *gin.Context) {

	companyId := ctx.DefaultQuery("company_id", "")

	if companyId == "" {
		utils.ErrorResponse(ctx, fmt.Errorf("company_id is required"), http.StatusBadRequest)
		return
	}

	// DefaultQuery returns the specified default value if the key is not found in the query string.
	pageNumber, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("size", "10"))

	offset := (pageNumber - 1) * pageSize

	// Query for get all employees associated to a company
	query := `
	SELECT 
		id, 
		name, 
		last_name, 
		phone_number, 
		email, 
		id_type, 
		id_number, 
		admission_date, 
		salary, 
		position_id, 
		department_id,
		company_id, 
		picture_url,
		created_at, 
		updated_at
	FROM employees
	WHERE company_id = $1
	ORDER BY id ASC
	LIMIT $2
	OFFSET $3
	`

	rows, err := s.db.Query(query, companyId, pageSize, offset)

	if err != nil {
		utils.ErrorResponse(ctx, err, http.StatusInternalServerError)
		return
	}

	// Query for get the total number of employees associated to a company
	totalItemsQuery := `SELECT COUNT(*) FROM employees WHERE company_id = $1`
	var totalItems int
	err = s.db.QueryRow(totalItemsQuery, companyId).Scan(&totalItems)

	if err != nil {
		utils.ErrorResponse(ctx, err, http.StatusInternalServerError)
		return
	}

	totalPages := int(math.Ceil(float64(totalItems) / float64(pageSize)))
	var nextPage, prevPage *int

	if pageNumber < totalPages {
		nextPageNum := pageNumber + 1
		nextPage = &nextPageNum
	}

	if pageNumber > 1 {
		prevPageNum := pageNumber - 1
		prevPage = &prevPageNum
	}

	var employees []models.EmployeeResponse

	for rows.Next() {
		var employee models.EmployeeResponse

		err := rows.Scan(
			&employee.ID,
			&employee.Name,
			&employee.LastName,
			&employee.PhoneNumber,
			&employee.Email,
			&employee.IdType,
			&employee.IdNumber,
			&employee.AdmissionDate,
			&employee.Salary,
			&employee.PositionId,
			&employee.DepartmentId,
			&employee.CompanyId,
			&employee.PictureUrl,
			&employee.CreatedAt,
			&employee.UpdatedAt)

		if err != nil {
			utils.ErrorResponse(ctx, err, http.StatusInternalServerError)
			return
		}

		employees = append(employees, employee)
	}

	result := models.PaginatedResult{
		Data:       employees,
		PageNumber: pageNumber,
		PageSize:   pageSize,
		TotalItems: totalItems,
		NextPage:   nextPage,
		PrevPage:   prevPage,
	}

	ctx.JSON(http.StatusOK, result)
}

func (s *Server) updateEmployee(ctx *gin.Context) {
	var params models.GetEmployeeParams
	var body models.CreateEmployeeBody

	if err := ctx.ShouldBindUri(&params); err != nil {
		utils.ErrorResponse(ctx, err, http.StatusBadRequest)
		return
	}

	if err := ctx.ShouldBindJSON(&body); err != nil {
		utils.ErrorResponse(ctx, err, http.StatusBadRequest)
		return
	}

	query := `
	UPDATE
		employees
	SET
		name = $1,
		last_name = $2,
		phone_number = $3,
		email = $4,
		id_type = $5,
		id_number = $6,
		admission_date = $7,
		salary = $8,
		position_id = $9,
		department_id = $10,
		company_id = $11
		picture_url = $12
	WHERE
		id = $13
	RETURNING id,
		name,
		last_name,
		phone_number,
		email,
		id_type,
		id_number,
		admission_date,
		salary,
		position_id,
		company_id,
		picture_url,
		created_at,
		updated_at
	`

	row := s.db.QueryRow(
		query,
		body.Name,
		body.LastName,
		body.PhoneNumber,
		body.Email,
		body.IdType,
		body.IdNumber,
		body.AdmissionDate,
		body.Salary,
		body.PositionId,
		body.CompanyId,
		body.PictureUrl,
		params.ID)

	var employee models.EmployeeResponse
	err := row.Scan(
		&employee.ID,
		&employee.Name,
		&employee.LastName,
		&employee.PhoneNumber,
		&employee.Email,
		&employee.IdType,
		&employee.IdNumber,
		&employee.AdmissionDate,
		&employee.Salary,
		&employee.PositionId,
		&employee.DepartmentId,
		&employee.CompanyId,
		&employee.PictureUrl,
		&employee.CreatedAt,
		&employee.UpdatedAt)

	if err != nil {
		utils.ErrorResponse(ctx, err, http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"employee": employee})
}

func (s *Server) deleteEmployee(ctx *gin.Context) {
	var params models.GetEmployeeParams

	if err := ctx.ShouldBindUri(&params); err != nil {
		utils.ErrorResponse(ctx, err, http.StatusBadRequest)
		return
	}

	query := `DELETE FROM employees WHERE id = $1`

	_, err := s.db.Exec(query, params.ID)

	if err != nil {
		utils.ErrorResponse(ctx, err, http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}
