package models

import "time"

type CreatePositionBody struct {
	Name         string `json:"name" binding:"required"`
	CompanyId    string `json:"company_id" binding:"required"`
	DepartmentId string `json:"department_id" binding:"required"`
}

type PositionResponse struct {
	ID           string
	Name         string
	CompanyId    string
	CreatedAt    time.Time
	DepartmentId string
	UpdatedAt    *time.Time
}

type GetCompanyPositionsParams struct {
	CompanyId string `uri:"companyId" binding:"required"`
}

type GetDepartmentPositionsParams struct {
	DepartmentId string `uri:"department_id" binding:"required"`
}

type GetPositionsParams struct {
	ID string `uri:"id" binding:"required"`
}

type SearchPositionsParams struct {
	CompanyId    string `form:"company_id"`
	DepartmentId string `form:"department_id"`
}
