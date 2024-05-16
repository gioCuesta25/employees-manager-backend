package models

type PaginatedResult struct {
	Data       any  `json:"data"`
	PageNumber int  `json:"pageNumber"`
	PageSize   int  `json:"pageSize"`
	TotalItems int  `json:"totalItems"`
	NextPage   *int `json:"nextPage"`
	PrevPage   *int `json:"prevPage"`
}
