package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gioCuesta25/employees-manager-backend/config"
	"github.com/golang-jwt/jwt"
)

type Server struct {
	env    config.Environment
	db     *sql.DB
	router *gin.Engine
}

func NewServer(env config.Environment, db *sql.DB) *Server {
	r := gin.Default()

	server := &Server{
		env:    env,
		db:     db,
		router: r,
	}

	// Routes
	server.setupRoutes()
	return server
}

func (s *Server) Run() {
	s.router.Run(s.env.ApiPort)
}

func (s *Server) setupRoutes() {

	// Users
	users := s.router.Group("/users")
	users.Use(s.RequireAuth)
	users.POST("/", s.createUser)
	users.GET("/", s.listUsers)
	users.GET("/:id", s.getUser)
	users.DELETE("/:id", s.deleteUser)
	users.PATCH("/:id", s.updateUser)

	//Login
	s.router.POST("/login", s.login)

	//Companies
	// Todo: Get companies by user
	companies := s.router.Group("/companies")
	companies.Use(s.RequireAuth)
	companies.POST("/", s.createCompany)
	companies.GET("/:id", s.getCompany)
	companies.DELETE("/:id", s.deleteCompany)
	companies.PATCH("/:id", s.updateCompany)

	// Departments
	departments := s.router.Group("/departments")
	departments.Use(s.RequireAuth)
	departments.POST("/", s.createDepartment)
	departments.GET("/:companyId", s.getDepartmentsByCompany)
	departments.PATCH("/:id", s.updateDepartment)
	departments.DELETE("/:id", s.deleteDepartment)
}

func (s *Server) RequireAuth(ctx *gin.Context) {
	// Get token from cookies or headers
	authorizationToken := ctx.GetHeader("Authorization")

	if authorizationToken == "" {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing authorization token"})
		return
	}
	tokenString := strings.Split(authorizationToken, "Bearer ")[1]

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.env.JwtSecret), nil
	})

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Check the exp
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "expired token"})
			return
		}

		// Find the user with token sub

		// Attach to req
		ctx.Set("userId", claims["sub"])

		// Continue
		ctx.Next()
	} else {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing authorization token"})
		return
	}
}
