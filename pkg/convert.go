package convert

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
}

// ToAPINote converts a note model to an API note
func ToAPINote(note *models.Note) *APINote {
	return &APINote{
		ID:           note.ID,
		UserID:       note.UserID,
		Title:        note.Title,
		Content:      note.Content,
		CurrentViews: note.CurrentViews,
	}
}

// ToAPINotes converts a slice of note models to a slice of API notes
func ToAPINotes(notes []*models.Note) []*APINote {
	apiNotes := make([]*APINote, 0)
	for _, note := range notes {
		apiNotes = append(apiNotes, ToAPINote(note))
	}
	return apiNotes
}
