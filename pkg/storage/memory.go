package storage

import (
	"fmt"
	"net/http"

	"github.com/codescalersinternships/secretnote-api-spa-eyadhussein/pkg/models"
	"github.com/codescalersinternships/secretnote-api-spa-eyadhussein/pkg/util"
)

// Memory is a storage implementation that uses memory
type Memory struct {
	notes          []*models.Note
	users          []*models.User
	usersIDCounter uint
	notesIDCounter uint
}

// NewMemory creates a new Memory storage
func NewMemory() *Memory {
	return &Memory{
		notes:          make([]*models.Note, 0),
		users:          make([]*models.User, 0),
		notesIDCounter: 1,
		usersIDCounter: 1,
	}
}

// GetUserByUsername gets a user by username
func (m *Memory) GetUserByUsername(username string) (*models.User, error) {
	for _, user := range m.users {
		if user.Username == username {
			return user, nil
		}
	}
	return nil, util.NewResponseError(fmt.Errorf("user with username %s not found", username), http.StatusNotFound)
}

// CreateUser creates a new user
func (m *Memory) CreateUser(user *models.User) error {
	for _, oldUser := range m.users {
		if oldUser.Username == user.Username {
			return util.NewResponseError(fmt.Errorf("user with username %s already exists", user.Username), http.StatusBadRequest)
		}
	}
	user.ID = m.usersIDCounter
	m.users = append(m.users, user)
	m.usersIDCounter++
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

// DeleteNoteByID deletes a note by ID
func (m *Memory) DeleteNoteByID(id string) error {
	return nil
}
