package models

import "time"

type CreateDepartmentBody struct {
	Name      string `json:"name" binding:"required"`
	CompanyId string `json:"company_id" binding:"required"`
}

type DepartmentsResponse struct {
	ID        string
	Name      string
	CompanyId string
	CreatedAt time.Time
	UpdatedAt *time.Time
}

type GetCompanyDepartmentsParams struct {
	CompanyId string `uri:"companyId" binding:"required"`
}

type GetDepartmentsParams struct {
	ID string `uri:"id" binding:"required"`
}
