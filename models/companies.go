package models

import "time"

type CreateCompanyBody struct {
	Name    string `json:"name" binding:"required"`
	Owner   int    `json:"owner" binding:"required"`
	Address string `json:"address" binding:"required"`
	Phone   string `json:"phone" binding:"required"`
	Email   string `json:"email" binding:"required"`
}

type GetCompanyParams struct {
	ID string `uri:"id" binding:"required"`
}

type CompanyResponse struct {
	ID        string     `json:"id"`
	Name      string     `json:"name"`
	Owner     int        `json:"owner"`
	Address   string     `json:"address"`
	Phone     string     `json:"phone"`
	Email     string     `json:"email"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}
