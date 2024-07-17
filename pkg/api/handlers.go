package api

import (
	"encoding/json"
	"net/http"

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

	user := models.NewUser(registerUserRequest.Username, registerUserRequest.Password, registerUserRequest.Email)

	err = s.store.CreateUser(user)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, nil)
}

func (s *Server) handleLoginUser(c *gin.Context) {
	loginUserRequest := &models.LoginUserRequest{}
	err := json.NewDecoder(c.Request.Body).Decode(loginUserRequest)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := s.store.GetUserByUsername(loginUserRequest.Username)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if user.Password != loginUserRequest.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid password"})
		return
	}

	c.JSON(http.StatusOK, nil)
}
