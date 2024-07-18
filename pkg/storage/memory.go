package storage

import (
	"github.com/codescalersinternships/secretnote-api-spa-eyadhussein/pkg/models"
)

// Memory is a storage implementation that uses memory
type Memory struct{}

// NewMemory creates a new Memory storage
func NewMemory() *Memory {
	return &Memory{}
}

// GetUserByUsername gets a user by username
func (m *Memory) GetUserByUsername(username string) (*models.User, error) {
	return nil, nil
}

// CreateUser creates a new user
func (m *Memory) CreateUser(user *models.User) error {
	return nil
}

// CreateNote creates a new note
func (m *Memory) CreateNote(note *models.Note) error {
	return nil
}

// GetNoteByID gets a note by ID
func (m *Memory) GetNoteByID(id string) (*models.Note, error) {
	return nil, nil
}

// GetNotesByUserID gets notes by user ID
func (m *Memory) GetNotesByUserID(userID int) ([]*models.Note, error) {
	return nil, nil
}

// IncrementNoteViews increments the views of a note
func (m *Memory) IncrementNoteViews(id string) error {
	return nil
}
