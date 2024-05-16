package models

import "time"

type CreateEmployeeBody struct {
	Name          string    `json:"name" binding:"required"`
	LastName      string    `json:"last_name" binding:"required"`
	PhoneNumber   string    `json:"phone_number" binding:"required"`
	Email         string    `json:"email" binding:"required"`
	IdType        string    `json:"id_type" binding:"required"`
	IdNumber      string    `json:"id_number" binding:"required"`
	AdmissionDate time.Time `json:"admission_date" binding:"required"`
	Salary        float64   `json:"salary" binding:"required"`
	PositionId    string    `json:"position_id" binding:"required"`
	DepartmentId  string    `json:"department_id" binding:"required"`
	CompanyId     string    `json:"company_id" binding:"required"`
}

type EmployeeResponse struct {
	ID            string
	Name          string
	LastName      string
	PhoneNumber   string
	Email         string
	IdType        string
	IdNumber      string
	AdmissionDate time.Time
	Salary        float64
	PositionId    string
	DepartmentId  string
	CompanyId     string
	CreatedAt     time.Time
	UpdatedAt     *time.Time
}

type GetEmployeeParams struct {
	ID string `uri:"id" binding:"required"`
}

type GetCompanyEmployeesParams struct {
	CompanyId string `uri:"companyId" binding:"required"`
}
