package models

import "gorm.io/gorm"

// User represents a user model
type User struct {
	gorm.Model
	ID       uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Username string `gorm:"size:100;unique;not null" json:"username"`
	Email    string `gorm:"size:100;unique;not null" json:"email"`
	Password string `gorm:"size:255;not null" json:"password"`
}

// NewUser creates a new user
func NewUser(username, email, password string) *User {
	return &User{
		Username: username,
		Email:    email,
		Password: password,
	}
}

// RegisterUserRequest represents a request to register a user
type RegisterUserRequest struct {
	Username             string `json:"username"`
	Email                string `json:"email"`
	Password             string `json:"password"`
	PasswordConfirmation string `json:"password_confirmation"`
}

// LoginUserRequest represents a request to login a user
type LoginUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
