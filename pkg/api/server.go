// Package api contains the server and routes
package api

import (
	"log"
	"time"

	"github.com/codescalersinternships/secretnote-api-spa-eyadhussein/pkg/api/middlewares"
	"github.com/codescalersinternships/secretnote-api-spa-eyadhussein/pkg/storage"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

const (
	rateLimit = 1
	burst     = 1
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
	config := cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	s.router.Use(cors.New(config))

	log.Printf("Starting server on port %s", s.listenAddr)

	err := s.router.Run(s.listenAddr)

	if err != nil {
		log.Fatalf("failed to start the server %v", err)
	}
}

func (s *Server) routes() {
	api := s.router.Group("/api")

	auth := api.Group("/auth")
	{
		auth.POST("/register", middlewares.RateLimiter(s.handleRegisterUser, rateLimit, burst))
		auth.POST("/login", middlewares.RateLimiter(s.handleLoginUser, rateLimit, burst))
		auth.POST("/logout", middlewares.VerifyToken(), middlewares.JwtAuthMiddleware(s.store), middlewares.RateLimiter(s.handleLogoutUser, rateLimit, burst))
		auth.POST("/verify-token", middlewares.VerifyToken())
	}

	api.GET("users/notes", middlewares.JwtAuthMiddleware(s.store), middlewares.RateLimiter(s.handleGetNotesByUserID, rateLimit, burst))

	notes := api.Group("/notes")
	{
		notes.POST("", middlewares.JwtAuthMiddleware(s.store), middlewares.RateLimiter(s.handleCreateNote, rateLimit, burst))
		notes.GET("/:id", middlewares.RateLimiter(s.handleGetNoteByID, rateLimit, burst))
	}

	api.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}
