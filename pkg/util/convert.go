package util

import (
	"github.com/codescalersinternships/secretnote-api-spa-eyadhussein/pkg/models"
	"github.com/google/uuid"
)

// APINote represents a note in API format
type APINote struct {
	ID           uuid.UUID `json:"id"`
	UserID       uint      `json:"user_id"`
	Title        string    `json:"title"`
	Content      string    `json:"content"`
	CurrentViews int       `json:"current_views"`
	MaxViews     int       `json:"max_views"`
	ExpiresAt    string    `json:"expires_at"`
}

// ToAPINote converts a note model to an API note
func ToAPINote(note *models.Note, restricted bool) *APINote {
	noteResponse := &APINote{
		ID:      note.ID,
		Title:   note.Title,
		Content: note.Content,
	}

	if !restricted {
		noteResponse.UserID = note.UserID
		noteResponse.MaxViews = note.MaxViews
		noteResponse.ExpiresAt = note.ExpiresAt.Format("2006-01-02T15:04:05Z07:00")
	}

	return noteResponse
}

// ToAPINotes converts a slice of note models to a slice of API notes
func ToAPINotes(notes []*models.Note, restricted bool) []*APINote {
	apiNotes := make([]*APINote, 0, len(notes))

	for _, note := range notes {
		apiNotes = append(apiNotes, ToAPINote(note, restricted))
	}

	return apiNotes
}
