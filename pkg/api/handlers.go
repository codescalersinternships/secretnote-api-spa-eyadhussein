package api

import (
	"net/http"
	"time"

	convert "github.com/codescalersinternships/secretnote-api-spa-eyadhussein/pkg"
	"github.com/codescalersinternships/secretnote-api-spa-eyadhussein/pkg/models"
	"github.com/gin-gonic/gin"
)

func (s *Server) handleRegisterUser(c *gin.Context) {
	var registerUserRequest models.RegisterUserRequest

	if err := c.ShouldBindJSON(&registerUserRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := models.NewUser(registerUserRequest.Username, registerUserRequest.Email, registerUserRequest.Password)

	err := user.SetPassword(user.Password)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := s.store.CreateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	token, err := createToken(registerUserRequest.Username, time.Hour*24*7)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie("token", token, 604800, "/", "localhost", true, true)
	c.SetCookie("user", user.Username, 604800, "/", "localhost", false, true)
	c.JSON(http.StatusCreated, gin.H{"message": "user registered successfully"})
}

func (s *Server) handleLoginUser(c *gin.Context) {
	var loginUserRequest models.LoginUserRequest

	if err := c.ShouldBindJSON(&loginUserRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := s.store.GetUserByUsername(loginUserRequest.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := user.CheckPassword(loginUserRequest.Password); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid password"})
		return
	}

	token, err := createToken(loginUserRequest.Username, time.Hour*24*7)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie("token", token, 604800, "/", "localhost", false, true)
	c.SetCookie("user", user.Username, 604800, "/", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "user login successfully"})
}

func (s *Server) handleLogoutUser(c *gin.Context) {
	c.SetCookie("token", "", 0, "/", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "successfully logged out"})
}

func (s *Server) handleCreateNote(c *gin.Context) {
	var createNoteRequest models.CreateNoteRequest

	if err := c.ShouldBindJSON(&createNoteRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "user not found"})
		return
	}

	authUser, ok := user.(*models.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user type"})
		return
	}

	note := models.NewNote(createNoteRequest.Title, createNoteRequest.Content, createNoteRequest.MaxViews, createNoteRequest.ExpiresAt)
	note.UserID = authUser.ID

	if err := s.store.CreateNote(note); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create note"})
		return
	}

	c.JSON(http.StatusCreated, convert.ToAPINote(note))
}

func (s *Server) handleGetNoteByID(c *gin.Context) {
	id := c.Param("id")

	note, err := s.store.GetNoteByID(id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not get note might have expired or reached max views"})
		return
	}

	if note.IsExpired() || note.HasReachedMaxViews() {
		if err := s.store.DeleteNoteByID(id); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not delete note"})
			return
		}

		c.JSON(http.StatusNotFound, gin.H{"error": "note has expired or reached max views"})
		return
	}

	if err := s.store.IncrementNoteViews(note.ID.String()); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not update note views"})
		return
	}

	note.CurrentViews++

	c.JSON(http.StatusOK, convert.ToAPINote(note))
}

func (s *Server) handleGetNotesByUserID(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "user not found"})
		return
	}

	authUser, ok := user.(*models.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user type"})
		return
	}

	notes, err := s.store.GetNotesByUserID(int(authUser.ID))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not get notes"})
		return
	}

	c.JSON(http.StatusOK, convert.ToAPINotes(notes))
}
