package api

import (
	"log"

	"github.com/codescalersinternships/secretnote-api-spa-eyadhussein/pkg/storage"
	"github.com/gin-gonic/gin"
)

// Server struct holds the listen address, storage and router
type Server struct {
	listenAddr string
	store      storage.Storage
	router     *gin.Engine
}

// NewServer creates a new server
func NewServer(listenAddr string, store storage.Storage) *Server {
	return &Server{listenAddr: listenAddr, store: store, router: gin.Default()}
}

// Run starts the server
func (s *Server) Run() {
	s.router.HandleMethodNotAllowed = true
	s.routes()

	log.Println("Starting server on port 8080...")

	err := s.router.Run(s.listenAddr)

	if err != nil {
		log.Fatalf("failed to start the server %v", err)
	}
}

func (s *Server) routes() {
	s.router.POST("/api/users/register", s.handleRegisterUser)
	s.router.POST("/api/users/login", s.handleLoginUser)
	s.router.POST("/api/users/logout", jwtAuthMiddlware, s.handleLogoutUser)
}
