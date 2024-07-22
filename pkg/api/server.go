// Package api contains the server and routes
package api

import (
	"log"
	"time"

	"github.com/codescalersinternships/secretnote-api-spa-eyadhussein/pkg/api/middlewares"
	"github.com/codescalersinternships/secretnote-api-spa-eyadhussein/pkg/storage"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Server struct holds the listen address, storage and router
type Server struct {
	listenAddr string
	store      storage.Storage
	router     *gin.Engine
	secretKey  string
	rateLimit  rate.Limit
	burst      int
}

// NewServer creates a new server
func NewServer(listenAddr string, store storage.Storage, secretKey string, rateLimit rate.Limit, burst int) *Server {
	if listenAddr == "" {
		listenAddr = ":5000"
	}
	return &Server{listenAddr: listenAddr, store: store, router: gin.Default(), secretKey: secretKey, rateLimit: rateLimit, burst: burst}
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

	err := s.router.Run("0.0.0.0" + s.listenAddr)

	if err != nil {
		log.Fatalf("failed to start the server %v", err)
	}
}

func (s *Server) routes() {
	api := s.router.Group("/api")

	auth := api.Group("/auth")
	{
		auth.POST("/register", middlewares.RateLimiter(s.handleRegisterUser, s.rateLimit, s.burst))
		auth.POST("/login", middlewares.RateLimiter(s.handleLoginUser, s.rateLimit, s.burst))
		auth.POST("/logout", middlewares.VerifyToken(s.secretKey), middlewares.RateLimiter(s.handleLogoutUser, s.rateLimit, s.burst))
		auth.POST("/verify-token", middlewares.VerifyToken(s.secretKey))
	}

	api.GET("users/notes", middlewares.JwtAuthMiddleware(s.store), middlewares.RateLimiter(s.handleGetNotesByUserID, s.rateLimit, s.burst))

	notes := api.Group("/notes")
	{
		notes.POST("", middlewares.JwtAuthMiddleware(s.store), middlewares.RateLimiter(s.handleCreateNote, s.rateLimit, s.burst))
		notes.GET("/:id", middlewares.RateLimiter(s.handleGetNoteByID, s.rateLimit, s.burst))
	}

	api.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}
