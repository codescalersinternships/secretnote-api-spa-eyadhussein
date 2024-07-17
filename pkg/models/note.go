package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Note represents a note model
type Note struct {
	gorm.Model
	ID           uuid.UUID `json:"id"`
	Title        string    `json:"title"`
	Content      string    `json:"content"`
	UserID       int       `json:"user_id"`
	MaxViews     int       `json:"max_views"`
	CurrentViews int       `json:"current_views"`
	ExpiresAt    time.Time `json:"expires_at"`
}

// IsExpired checks if the note is expired based on the expiration date
func (n *Note) IsExpired() bool {
	return time.Now().After(n.ExpiresAt)
}

// HasReachedMaxViews checks if the note has reached the maximum allowed views
func (n *Note) HasReachedMaxViews() bool {
	return n.CurrentViews >= n.MaxViews
}
