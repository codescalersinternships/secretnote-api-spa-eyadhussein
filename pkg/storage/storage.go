// Storage package defines the Storage interface that defines the methods that a storage implementation must have.
package storage

import (
	"github.com/codescalersinternships/secretnote-api-spa-eyadhussein/pkg/models"
)

// Storage is an interface that defines the methods that a storage implementation must have
type Storage interface {
	CreateUser(user *models.User) error
	GetUserByUsername(username string) (*models.User, error)
	CreateNote(note *models.Note) error
	GetNoteByID(id string) (*models.Note, error)
	GetNotesByUserID(userID int) ([]*models.Note, error)
	IncrementNoteViews(id string) error
}
