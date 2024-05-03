package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gioCuesta25/employees-manager-backend/models"
	"github.com/gioCuesta25/employees-manager-backend/utils"
)

func (s *Server) createUser(ctx *gin.Context) {
	var body models.CreateUserRequest

	if err := ctx.ShouldBindJSON(&body); err != nil {
		utils.ErrorResponse(ctx, err, http.StatusBadRequest)
		return
	}

	hashedPassword, err := utils.HashPassword(body.Password)

	if err != nil {
		utils.ErrorResponse(ctx, err, http.StatusInternalServerError)
		return
	}

	query := "INSERT INTO users (full_name, email, password) VALUES ($1, $2, $3) RETURNING id, full_name, email"

	row := s.db.QueryRow(query, body.FullName, body.Email, hashedPassword)

	var user models.GetUsersResponse

	err = row.Scan(&user.ID, &user.FullName, &user.Email)

	if err != nil {
		utils.ErrorResponse(ctx, err, http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"user": user})
}

func (s *Server) listUsers(ctx *gin.Context) {
	page := 1
	size := 25

	//TODO: Improve pagination
	//Todo: Add created_at and updated_at in the response

	if value, ok := ctx.GetQuery("page"); ok {
		conv, err := strconv.Atoi(value)
		if err != nil {
			utils.ErrorResponse(ctx, err, http.StatusInternalServerError)
			return
		}
		page = conv
	}

	if value, ok := ctx.GetQuery("size"); ok {
		conv, err := strconv.Atoi(value)
		if err != nil {
			utils.ErrorResponse(ctx, err, http.StatusInternalServerError)
			return
		}
		size = conv
	}

	offset := (page - 1) * size

	rows, err := s.db.Query("SELECT id, full_name, email FROM users LIMIT $1 OFFSET $2", size, offset)

	if err != nil {
		utils.ErrorResponse(ctx, err, http.StatusInternalServerError)
	}

	users := make([]*models.GetUsersResponse, 0)

	for rows.Next() {
		u, err := scanRowsIntoUser(rows)
		if err != nil {
			utils.ErrorResponse(ctx, err, http.StatusInternalServerError)
			return
		}
		users = append(users, u)
	}

	ctx.JSON(http.StatusOK, gin.H{"items": users})
}

func (s *Server) getUser(ctx *gin.Context) {
	var params models.GetUserRequest
	if err := ctx.ShouldBindUri(&params); err != nil {
		utils.ErrorResponse(ctx, err, http.StatusBadRequest)
		return
	}

	query := "SELECT id, full_name, email FROM users WHERE id = $1"

	row := s.db.QueryRow(query, params.ID)

	var user models.GetUsersResponse

	err := row.Scan(&user.ID, &user.FullName, &user.Email)

	if err != nil {
		if err == sql.ErrNoRows {
			utils.ErrorResponse(ctx, fmt.Errorf("not found user with id %s", params.ID), http.StatusNotFound)
			return
		}
		utils.ErrorResponse(ctx, err, http.StatusBadRequest)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"user": user})
}

func (s *Server) deleteUser(ctx *gin.Context) {
	var params models.GetUserRequest

	if err := ctx.ShouldBindUri(&params); err != nil {
		utils.ErrorResponse(ctx, err, http.StatusBadRequest)
		return
	}

	query := "DELETE FROM users WHERE id = $1"

	_, err := s.db.Query(query, params.ID)

	if err != nil {
		utils.ErrorResponse(ctx, err, http.StatusBadRequest)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})

}

func (s *Server) updateUser(ctx *gin.Context) {
	var body models.CreateUserResponse
	var params models.GetUserRequest

	if err := ctx.ShouldBindJSON(&body); err != nil {
		utils.ErrorResponse(ctx, err, http.StatusBadRequest)
		return
	}

	if err := ctx.ShouldBindUri(&params); err != nil {
		utils.ErrorResponse(ctx, err, http.StatusBadRequest)
		return
	}

	query := `UPDATE users
			SET full_name = $1, email = $2
			WHERE id = $3
			RETURNING id, full_name, email`

	row := s.db.QueryRow(query, body.FullName, body.Email, params.ID)

	var user models.GetUsersResponse

	err := row.Scan(&user.ID, &user.FullName, &user.Email)

	if err != nil {
		if err == sql.ErrNoRows {
			utils.ErrorResponse(ctx, fmt.Errorf("not found user with id %s", params.ID), http.StatusNotFound)
			return
		}
		utils.ErrorResponse(ctx, err, http.StatusBadRequest)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"user": user})
}

func (s *Server) login(ctx *gin.Context) {
	var body models.LoginBody

	if err := ctx.ShouldBindJSON(&body); err != nil {
		utils.ErrorResponse(ctx, err, http.StatusBadRequest)
		return
	}

	query := "SELECT * FROM users WHERE email = $1"

	row := s.db.QueryRow(query, body.Email)

	var user models.CompleteUserResponse

	err := row.Scan(&user.ID, &user.FullName, &user.Email, &user.Password)

	if err != nil {
		if err == sql.ErrNoRows {
			utils.ErrorResponse(ctx, fmt.Errorf("not found with email %s", body.Email), http.StatusNotFound)
			return
		}
		utils.ErrorResponse(ctx, err, http.StatusBadRequest)
		return
	}

	matchPassword := utils.CheckPasswordHash(body.Password, user.Password)

	if !matchPassword {
		utils.ErrorResponse(ctx, fmt.Errorf("invalid credentials"), http.StatusUnauthorized)
		return
	}

	//Generate jwt token
	token, err := utils.GetToken(user.ID, s.env.JwtSecret)

	if err != nil {
		utils.ErrorResponse(ctx, err, http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}

func scanRowsIntoUser(rows *sql.Rows) (*models.GetUsersResponse, error) {
	usser := new(models.GetUsersResponse)

	err := rows.Scan(
		&usser.ID,
		&usser.FullName,
		&usser.Email,
	)

	if err != nil {
		return nil, err
	}

	return usser, nil
}
