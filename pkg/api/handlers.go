package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/codescalersinternships/secretnote-api-spa-eyadhussein/pkg/models"
	"github.com/gin-gonic/gin"
)

func (s *Server) handleRegisterUser(c *gin.Context) {
	registerUserRequest := &models.RegisterUserRequest{}
	err := json.NewDecoder(c.Request.Body).Decode(registerUserRequest)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := registerUserRequest.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := models.NewUser(registerUserRequest.Username, registerUserRequest.Email, registerUserRequest.Password)

	err = user.SetPassword(user.Password)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = s.store.CreateUser(user)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	token, err := createToken(registerUserRequest.Username, time.Hour*24*7)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"token": token})
}

func (s *Server) handleLoginUser(c *gin.Context) {
	loginUserRequest := &models.LoginUserRequest{}
	err := json.NewDecoder(c.Request.Body).Decode(loginUserRequest)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := loginUserRequest.Validate(); err != nil {
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

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (s *Server) handleLogoutUser(c *gin.Context) {
	c.Set("Authorization", "")
	c.JSON(http.StatusOK, gin.H{"message": "successfully logged out"})
}
