package api

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gioCuesta25/employees-manager-backend/models"
	"github.com/gioCuesta25/employees-manager-backend/utils"
)

func (s *Server) signUp(ctx *gin.Context) {
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

	//Todo: Implement logic for token generation

	ctx.JSON(http.StatusOK, gin.H{
		"user": &models.CreateUserResponse{
			FullName: user.FullName,
			Email:    user.Email,
		},
		"token": "aquí va el token",
	})
}
