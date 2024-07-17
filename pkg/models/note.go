package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Note represents a note model
type Note struct {
	gorm.Model
	ID           uuid.UUID `gorm:"varchar(130);primaryKey" json:"id"`
	Title        string    `gorm:"type:varchar(255)" json:"title"`
	Content      string    `gorm:"type:text" json:"content"`
	UserID       uint      `gorm:"not null" json:"user_id"`
	User         *User     `gorm:"foreignKey:UserID" json:"user"`
	MaxViews     int       `json:"max_views"`
	CurrentViews int       `json:"current_views"`
	ExpiresAt    time.Time `json:"expires_at"`
}

// NewNote creates a new note
func NewNote(title, content string, maxViews int, expiresAt time.Time, userID uint) *Note {
	return &Note{
		Title:     title,
		Content:   content,
		MaxViews:  maxViews,
		ExpiresAt: expiresAt,
		UserID:    userID,
	}
}

// BeforeCreate is a GORM hook that is triggered before a new record is inserted into the database
func (note *Note) BeforeCreate(tx *gorm.DB) (err error) {
	note.ID = uuid.New()
	return
}

// IsExpired checks if the note is expired based on the expiration date
func (n *Note) IsExpired() bool {
	return time.Now().After(n.ExpiresAt)
}

// HasReachedMaxViews checks if the note has reached the maximum allowed views
func (n *Note) HasReachedMaxViews() bool {
	return n.CurrentViews >= n.MaxViews
}
