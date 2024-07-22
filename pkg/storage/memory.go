package storage

import (
	"strconv"

	"github.com/codescalersinternships/secretnote-api-spa-eyadhussein/pkg/models"
	"github.com/codescalersinternships/secretnote-api-spa-eyadhussein/pkg/util"
	"github.com/google/uuid"
)

// Memory is a storage implementation that uses memory
type Memory struct {
	notes          []*models.Note
	users          []*models.User
	UsersIDCounter uint
	NotesIDCounter uint
}

// NewMemory creates a new Memory storage
func NewMemory() *Memory {
	return &Memory{
		notes:          make([]*models.Note, 0),
		users:          make([]*models.User, 0),
		UsersIDCounter: 1,
		NotesIDCounter: 1000000000000000,
	}
}

// GetUserByUsername gets a user by username
func (m *Memory) GetUserByUsername(username string) (*models.User, error) {
	for _, user := range m.users {
		if user.Username == username {
			return user, nil
		}
	}
	return nil, util.ErrNotFound
}

// CreateUser creates a new user
func (m *Memory) CreateUser(user *models.User) error {
	for _, oldUser := range m.users {
		if oldUser.Username == user.Username {
			return util.ErrBadRequest
		}
	}
	user.ID = m.UsersIDCounter
	m.users = append(m.users, user)
	m.UsersIDCounter++
	return nil
}

// CreateNote creates a new note
func (m *Memory) CreateNote(note *models.Note) error {
	noteID := m.NotesIDCounter
	note.ID = uuid.UUID([]byte(strconv.Itoa(int(noteID))))
	m.notes = append(m.notes, note)
	m.NotesIDCounter++
	return nil
}

// GetNoteByID gets a note by ID
func (m *Memory) GetNoteByID(id string) (*models.Note, error) {
	for _, note := range m.notes {
		if note.ID.String() == id {
			return note, nil
		}
	}
	return nil, util.ErrNotFound
}

// GetNotesByUserID gets notes by user ID
func (m *Memory) GetNotesByUserID(userID uint) ([]*models.Note, error) {
	notes := make([]*models.Note, 0)
	for _, note := range m.notes {
		if note.UserID == userID {
			notes = append(notes, note)
		}
	}
	return notes, nil
}

// IncrementNoteViews increments the views of a note
func (m *Memory) IncrementNoteViews(id string) error {
	_, err := m.GetNoteByID(id)
	if err != nil {
		return err
	}
	return nil
}

// DeleteNoteByID deletes a note by ID
func (m *Memory) DeleteNoteByID(id string) error {
	for i, note := range m.notes {
		if note.ID.String() == id {
			m.notes = append(m.notes[:i], m.notes[i+1:]...)
			return nil
		}
	}
	return util.ErrNotFound
}

// Clear clears the storage
func (m *Memory) Clear() error {
	m.notes = make([]*models.Note, 0)
	m.users = make([]*models.User, 0)
	m.UsersIDCounter = 1
	m.NotesIDCounter = 1000000000000000
	return nil
}
