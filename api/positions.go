package api

import (
	"database/sql"
	"net/http"
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

func (s *Server) getPositionsByDepartment(ctx *gin.Context) {
	var params models.GetDepartmentPositionsParams

	if err := ctx.ShouldBindUri(&params); err != nil {
		utils.ErrorResponse(ctx, err, http.StatusBadRequest)
		return
	}

	query := `SELECT id, name, company_id, department_id, created_at, updated_at FROM positions WHERE company_id = $1`

	rows, err := s.db.Query(query, params.DepartmentId)

	positions := make([]*models.PositionResponse, 0)

	for rows.Next() {
		p, err := scanRowsIntoPosition(rows)
		if err != nil {
			utils.ErrorResponse(ctx, err, http.StatusInternalServerError)
			return
		}
		positions = append(positions, p)
	}

	if err != nil {
		utils.ErrorResponse(ctx, err, http.StatusInternalServerError)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"items": positions})
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
