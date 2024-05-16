package api

import (
	"net/http"

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
		company_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, name, last_name, phone_number, email, id_type, id_number, admission_date, salary, position_id, company_id, created_at, updated_at`

	row := s.db.QueryRow(query, body.Name, body.LastName, body.PhoneNumber, body.Email, body.IdType, body.IdNumber, body.AdmissionDate, body.Salary, body.PositionId, body.CompanyId)

	var employee models.EmployeeResponse

	err := row.Scan(&employee.ID, &employee.Name, &employee.LastName, &employee.PhoneNumber, &employee.Email, &employee.IdType, &employee.IdNumber, &employee.AdmissionDate, &employee.Salary, &employee.PositionId, &employee.CompanyId, &employee.CreatedAt, &employee.UpdatedAt)

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

	query := `SELECT id, name, last_name, phone_number, email, id_type, id_number, admission_date, salary, position_id, company_id, created_at, updated_at FROM employees WHERE id = $1`

	row := s.db.QueryRow(query, params.ID)

	var employee models.EmployeeResponse

	err := row.Scan(&employee.ID, &employee.Name, &employee.LastName, &employee.PhoneNumber, &employee.Email, &employee.IdType, &employee.IdNumber, &employee.AdmissionDate, &employee.Salary, &employee.PositionId, &employee.CompanyId, &employee.CreatedAt, &employee.UpdatedAt)

	if err != nil {
		utils.ErrorResponse(ctx, err, http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"employee": employee})
}

func (s *Server) listCompanyEmployees(ctx *gin.Context) {

	//Todo: Implement pagination
	var params models.GetCompanyEmployeesParams

	if err := ctx.ShouldBindUri(&params); err != nil {
		utils.ErrorResponse(ctx, err, http.StatusBadRequest)
		return
	}

	query := `SELECT id, name, last_name, phone_number, email, id_type, id_number, admission_date, salary, position_id, company_id, created_at, updated_at FROM employees WHERE company_id = $1`

	rows, err := s.db.Query(query, params.CompanyId)

	if err != nil {
		utils.ErrorResponse(ctx, err, http.StatusInternalServerError)
		return
	}

	var employees []models.EmployeeResponse

	for rows.Next() {
		var employee models.EmployeeResponse

		err := rows.Scan(&employee.ID, &employee.Name, &employee.LastName, &employee.PhoneNumber, &employee.Email, &employee.IdType, &employee.IdNumber, &employee.AdmissionDate, &employee.Salary, &employee.PositionId, &employee.CompanyId, &employee.CreatedAt, &employee.UpdatedAt)

		if err != nil {
			utils.ErrorResponse(ctx, err, http.StatusInternalServerError)
			return
		}

		employees = append(employees, employee)
	}

	ctx.JSON(http.StatusOK, gin.H{"employees": employees})
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

	query := `UPDATE employees SET name = $1, last_name = $2, phone_number = $3, email = $4, id_type = $5, id_number = $6, admission_date = $7, salary = $8, position_id = $9, company_id = $10 WHERE id = $11 RETURNING id, name, last_name, phone_number, email, id_type, id_number, admission_date, salary, position_id, company_id, created_at, updated_at`

	row := s.db.QueryRow(query, body.Name, body.LastName, body.PhoneNumber, body.Email, body.IdType, body.IdNumber, body.AdmissionDate, body.Salary, body.PositionId, body.CompanyId, params.ID)

	var employee models.EmployeeResponse
	err := row.Scan(&employee.ID, &employee.Name, &employee.LastName, &employee.PhoneNumber, &employee.Email, &employee.IdType, &employee.IdNumber, &employee.AdmissionDate, &employee.Salary, &employee.PositionId, &employee.CompanyId, &employee.CreatedAt, &employee.UpdatedAt)

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
