package api

import (
	"net/http"
	"time"

	"github.com/codescalersinternships/secretnote-api-spa-eyadhussein/pkg/api/middlewares"
	"github.com/codescalersinternships/secretnote-api-spa-eyadhussein/pkg/models"
	"github.com/codescalersinternships/secretnote-api-spa-eyadhussein/pkg/util"
	"github.com/gin-gonic/gin"
)

// @Summary Register a new user
// @Description Register a new user with username, email, and password
// @Tags auth
// @Accept json
// @Produce json
// @Param body body models.RegisterUserRequest true "User credentials to register"
// @Success 201 {object} swagger.ResponseUserRegistered
// @Failure 400 {object} swagger.ResponseBadRequest
// @Failure 500 {object} swagger.ResponseInternalServerError
// @Router /auth/register [post]
func (s *Server) handleRegisterUser(c *gin.Context) {
	var registerUserRequest models.RegisterUserRequest

	if err := c.ShouldBindJSON(&registerUserRequest); err != nil {
		c.JSON(http.StatusBadRequest, util.NewResponseError(util.ErrBadRequest, http.StatusBadRequest))
		return
	}

	user := models.NewUser(registerUserRequest.Username, registerUserRequest.Email, registerUserRequest.Password)

	err := user.SetPassword(user.Password)

	if err != nil {
		c.JSON(http.StatusInternalServerError, util.NewResponseError(util.ErrInternalServer, http.StatusInternalServerError))
		return
	}

	if err := s.store.CreateUser(user); err != nil {
		c.JSON(http.StatusBadRequest, util.NewResponseError(err, http.StatusBadRequest))
		return
	}

	token, err := middlewares.CreateToken(registerUserRequest.Username, time.Hour*24*7)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.NewResponseError(err, http.StatusInternalServerError))
		return
	}

	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie("token", token, 604800, "/", "localhost", true, true)
	c.SetCookie("user", user.Username, 604800, "/", "localhost", false, true)
	c.JSON(http.StatusCreated, gin.H{"message": "user registered successfully"})
}

// @Summary Login a user
// @Description Login a user
// @Tags auth
// @Accept json
// @Produce json
// @Param body body models.LoginUserRequest true "User credentials to login"
// @Success 200 {object} swagger.ResponseUserLoggedIn
// @Failure 400 {object} swagger.ResponseBadRequest
// @Failure 401 {object} swagger.ResponseUnauthorized
// @Failure 404 {object} swagger.ResponseNotFound
// @Failure 500 {object} swagger.ResponseInternalServerError
// @Router /auth/login [post]
func (s *Server) handleLoginUser(c *gin.Context) {
	var loginUserRequest models.LoginUserRequest

	if err := c.ShouldBindJSON(&loginUserRequest); err != nil {
		c.JSON(http.StatusBadRequest, util.NewResponseError(util.ErrBadRequest, http.StatusBadRequest))
		return
	}

	user, err := s.store.GetUserByUsername(loginUserRequest.Username)
	if err != nil {
		c.JSON(http.StatusNotFound, util.NewResponseError(err, http.StatusNotFound))
		return
	}

	if err := user.CheckPassword(loginUserRequest.Password); err != nil {
		c.JSON(http.StatusUnauthorized, util.NewResponseError(util.ErrUnauthorized, http.StatusUnauthorized))
		return
	}

	token, err := middlewares.CreateToken(loginUserRequest.Username, time.Hour*24*7)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.NewResponseError(err, http.StatusInternalServerError))
		return
	}

	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie("token", token, 604800, "/", "localhost", false, true)
	c.SetCookie("user", user.Username, 604800, "/", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "user login successfully"})
}

// @Summary Logout a user
// @Description Logout a user
// @Tags auth
// @Produce json
// @Success 200 {object} swagger.ResponseUserLoggedOut
// @Failure 401 {object} swagger.ResponseUnauthorized
// @Failure 500 {object} swagger.ResponseInternalServerError
// @Router /auth/logout [post]
// @Security Token
func (s *Server) handleLogoutUser(c *gin.Context) {
	c.SetCookie("token", "", 0, "/", "localhost", false, true)
	c.SetCookie("user", "", 0, "/", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "successfully logged out"})
}

// @Summary Create a note
// @Description Create a note with title, content, max views, and expiration date
// @Tags notes
// @Accept json
// @Produce json
// @Param body body models.CreateNoteRequest true "Note details to create"
// @Success 201 {object} util.APINote
// @Failure 400 {object} swagger.ResponseBadRequest
// @Failure 401 {object} swagger.ResponseUnauthorized
// @Failure 404 {object} swagger.ResponseNotFound
// @Failure 500 {object} swagger.ResponseInternalServerError
// @Router /notes [post]
// @Security Token
func (s *Server) handleCreateNote(c *gin.Context) {
	var createNoteRequest models.CreateNoteRequest

	if err := c.ShouldBindJSON(&createNoteRequest); err != nil {
		c.JSON(http.StatusBadRequest, util.NewResponseError(util.ErrBadRequest, http.StatusBadRequest))
		return
	}

	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusNotFound, util.NewResponseError(util.ErrNotFound, http.StatusNotFound))
		return
	}

	authUser, ok := user.(*models.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, util.NewResponseError(util.ErrInternalServer, http.StatusInternalServerError))
		return
	}

	note := models.NewNote(createNoteRequest.Title, createNoteRequest.Content, createNoteRequest.MaxViews, createNoteRequest.ExpiresAt)
	note.UserID = authUser.ID

	if err := s.store.CreateNote(note); err != nil {
		c.JSON(http.StatusInternalServerError, util.NewResponseError(err, http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusCreated, util.ToAPINote(note, false))
}

// @Summary Get a note by ID
// @Description Get a note by ID
// @Tags notes
// @Accept json
// @Produce json
// @Param id path string true "Note ID"
// @Success 200 {object} swagger.ResponseNoteRetrievedRestricted
// @Failure 404 {object} swagger.ResponseNotFound
// @Failure 500 {object} swagger.ResponseInternalServerError
// @Router /notes/{id} [get]
func (s *Server) handleGetNoteByID(c *gin.Context) {
	id := c.Param("id")

	note, err := s.store.GetNoteByID(id)

	if err != nil {
		c.JSON(http.StatusNotFound, util.NewResponseError(err, http.StatusNotFound))
		return
	}

	if note.IsExpired() || note.HasReachedMaxViews() {
		if err := s.store.DeleteNoteByID(id); err != nil {
			c.JSON(http.StatusInternalServerError, util.NewResponseError(err, http.StatusInternalServerError))
			return
		}

		c.JSON(http.StatusNotFound, util.NewResponseError(err, http.StatusNotFound))
		return
	}

	if err := s.store.IncrementNoteViews(note.ID.String()); err != nil {
		c.JSON(http.StatusInternalServerError, util.NewResponseError(err, http.StatusInternalServerError))
		return
	}

	note.CurrentViews++

	c.JSON(http.StatusOK, util.ToAPINote(note, true))
}

// @Summary Get notes by user ID
// @Description Get notes by user ID
// @Tags notes
// @Produce json
// @Success 200 {object} []util.APINote
// @Failure 401 {object} swagger.ResponseUnauthorized
// @Failure 404 {object} swagger.ResponseNotFound
// @Failure 500 {object} swagger.ResponseInternalServerError
// @Router /users/notes [get]
// @Security Token
func (s *Server) handleGetNotesByUserID(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusNotFound, util.NewResponseError(util.ErrNotFound, http.StatusNotFound))
		return
	}

	authUser, ok := user.(*models.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, util.NewResponseError(util.ErrInternalServer, http.StatusInternalServerError))
		return
	}

	notes, err := s.store.GetNotesByUserID(int(authUser.ID))

	if err != nil {
		c.JSON(http.StatusNotFound, util.NewResponseError(err, http.StatusNotFound))
		return
	}

	c.JSON(http.StatusOK, util.ToAPINotes(notes, false))
}
