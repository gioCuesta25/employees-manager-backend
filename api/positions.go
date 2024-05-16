package api

import (
	"database/sql"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gioCuesta25/employees-manager-backend/models"
	"github.com/gioCuesta25/employees-manager-backend/utils"
)

func (s *Server) createPosition(ctx *gin.Context) {
	var body models.CreatePositionBody

	if err := ctx.ShouldBindJSON(&body); err != nil {
		utils.ErrorResponse(ctx, err, http.StatusBadRequest)
		return
	}

	query := `INSERT INTO position
	(name, company_id, department_id)
	VALUES ($1, $2, $3)
	RETURNING id, name, company_id, department_id, created_at, updated_at`

	row := s.db.QueryRow(query, body.Name, body.CompanyId, body.DepartmentId)

	var department models.PositionResponse

	err := row.Scan(&department.ID,
		&department.Name,
		&department.CompanyId,
		&department.DepartmentId,
		&department.CreatedAt,
		&department.UpdatedAt)

	if err != nil {
		utils.ErrorResponse(ctx, err, http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"department": department})
}

func (s *Server) searchPositions(ctx *gin.Context) {
	departmentId := ctx.DefaultQuery("department_id", "")
	companyId := ctx.DefaultQuery("company_id", "")

	pageNumber, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("size", "10"))

	if departmentId == "" && companyId == "" {
		utils.ErrorResponse(ctx, fmt.Errorf("at least one param is required"), http.StatusBadRequest)
		return
	}

	query := `SELECT id, name, company_id, department_id, created_at, updated_at
	FROM positions
	WHERE ($1 = '' OR department_id = $1) AND ($2 = '' OR company_id = $2)
	ORDER BY id
	LIMIT $3
	OFFSET $4`

	offset := (pageNumber - 1) * pageSize

	rows, err := s.db.Query(query, departmentId, companyId, pageSize, offset)

	if err != nil {
		utils.ErrorResponse(ctx, err, http.StatusInternalServerError)
		return

	}

	totalItemsQuery := `SELECT COUNT(*) FROM positions WHERE ($1 = '' OR department_id = $1) AND ($2 = '' OR company_id = $2)`
	var totalItems int
	err = s.db.QueryRow(totalItemsQuery, departmentId, companyId).Scan(&totalItems)

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
	var positions []models.PositionResponse

	for rows.Next() {
		position, err := scanRowsIntoPosition(rows)

		if err != nil {
			utils.ErrorResponse(ctx, err, http.StatusInternalServerError)
			return
		}

		positions = append(positions, *position)
	}

	result := models.PaginatedResult{
		Data:       positions,
		PageNumber: pageNumber,
		PageSize:   pageSize,
		TotalItems: totalItems,
		NextPage:   nextPage,
		PrevPage:   prevPage,
	}

	ctx.JSON(http.StatusOK, result)
}

func (s *Server) updatePosition(ctx *gin.Context) {
	var body models.CreatePositionBody
	var params models.GetPositionsParams

	if err := ctx.ShouldBindJSON(&body); err != nil {
		utils.ErrorResponse(ctx, err, http.StatusBadRequest)
		return
	}

	if err := ctx.ShouldBindUri(&params); err != nil {
		utils.ErrorResponse(ctx, err, http.StatusBadRequest)
		return
	}

	query := `UPDATE positions
	SET name = $1, updated_at = $2
	WHERE id = $3
	RETURNING id, name, company_id, department_id, created_at, updated_at`

	row := s.db.QueryRow(query, body.Name, body.CompanyId, time.Now(), params.ID)

	var department models.PositionResponse

	err := row.Scan(&department.ID,
		&department.Name,
		&department.CompanyId,
		&department.DepartmentId,
		&department.CreatedAt,
		&department.UpdatedAt)

	if err != nil {
		utils.ErrorResponse(ctx, err, http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"department": department})
}

func (s *Server) deletePosition(ctx *gin.Context) {
	var params models.GetPositionsParams

	if err := ctx.ShouldBindUri(&params); err != nil {
		utils.ErrorResponse(ctx, err, http.StatusBadRequest)
		return
	}

	query := `DELETE FROM positions WHERE id = $1`

	_, err := s.db.Query(query, params.ID)

	if err != nil {
		utils.ErrorResponse(ctx, err, http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Department deleted successfully"})

}

func scanRowsIntoPosition(rows *sql.Rows) (*models.PositionResponse, error) {
	position := new(models.PositionResponse)

	err := rows.Scan(
		&position.ID,
		&position.Name,
		&position.CompanyId,
		&position.DepartmentId,
		&position.CreatedAt,
		&position.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return position, nil
}
