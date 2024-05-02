package api

import (
	"database/sql"

	"github.com/gin-gonic/gin"
)

type Server struct {
	listedAddr string
	db         *sql.DB
	router     *gin.Engine
}

func NewServer(lisenAddr string, db *sql.DB) *Server {
	r := gin.Default()

	server := &Server{
		listedAddr: lisenAddr,
		db:         db,
		router:     r,
	}

	// Routes
	server.setupRoutes()
	return server
}

func (s *Server) Run() {
	s.router.Run(s.listedAddr)
}

func (s *Server) setupRoutes() {

	// Users
	s.router.POST("/users", s.createUser)
	s.router.GET("/users", s.listUsers)
	s.router.GET("/users/:id", s.getUser)
	s.router.DELETE("/users/:id", s.deleteUser)
	s.router.PATCH("/users/:id", s.updateUser)

	//Login
	s.router.POST("/login", s.login)
}
