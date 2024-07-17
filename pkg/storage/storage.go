package storage

import "github.com/codescalersinternships/secretnote-api-spa-eyadhussein/pkg/models"

type Storage interface {
	CreateUser(user *models.User) error
	GetUserByUsername(username string) (*models.User, error)
}
