package models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// User represents a user model
type User struct {
	gorm.Model
	ID       uint    `gorm:"primaryKey;autoIncrement" json:"id"`
	Username string  `gorm:"size:100;unique;not null" json:"username"`
	Email    string  `gorm:"size:100;unique;not null" json:"email"`
	Password string  `gorm:"size:255;not null" json:"password"`
	Notes    []*Note `gorm:"foreignKey:UserID" json:"notes"`
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
// @Description RegisterUserRequest represents a request to register a user
// @Model
type RegisterUserRequest struct {
	// @Description The username of the user
	// @Example johndoe
	Username string `json:"username" binding:"required"`

	// @Description The email address of the user
	// @Example johndoe@example.com
	Email string `json:"email" binding:"required,email"`

	// @Description The password for the user
	// @Example P@ssw0rd!
	Password string `json:"password" binding:"required,min=8"`

	// @Description Confirmation of the password
	// @Example P@ssw0rd!
	PasswordConfirmation string `json:"password_confirmation" binding:"required,eqfield=Password"`
}

// LoginUserRequest represents a request to login a user
type LoginUserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=8"`
}

// SetPassword hashes the password using bcrypt
func (u *User) SetPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// CheckPassword checks if the provided password matches the stored hashed password
func (u *User) CheckPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}
