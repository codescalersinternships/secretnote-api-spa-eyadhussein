// Package api contains the server and routes
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

	log.Printf("Starting server on port %s", s.listenAddr)

	err := s.router.Run(s.listenAddr)

	if err != nil {
		log.Fatalf("failed to start the server %v", err)
	}
}

func (s *Server) routes() {
	api := s.router.Group("/api")

	users := api.Group("/users")
	{
		users.POST("/register", s.handleRegisterUser)
		users.POST("/login", s.handleLoginUser)
		users.POST("/logout", jwtAuthMiddleware(s.store), s.handleLogoutUser)
		users.GET("/notes", jwtAuthMiddleware(s.store), s.handleGetNotesByUserID)
	}

	notes := api.Group("/notes")
	{
		notes.POST("", jwtAuthMiddleware(s.store), s.handleCreateNote)
		notes.GET("/:id", jwtAuthMiddleware(s.store), s.handleGetNoteByID)
	}
}
