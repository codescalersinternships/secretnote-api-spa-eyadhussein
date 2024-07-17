package models

import (
	"errors"
	"regexp"

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

// Validate checks if the LoginUserRequest is valid
func (r *LoginUserRequest) Validate() error {
	if r.Username == "" {
		return errors.New("username is required")
	}
	if r.Password == "" {
		return errors.New("password is required")
	}
	return nil
}

// Validate validates the RegisterUserRequest fields
func (r *RegisterUserRequest) Validate() error {
	if r.Username == "" {
		return errors.New("username is required")
	}
	if r.Email == "" {
		return errors.New("email is required")
	}
	if !isValidEmail(r.Email) {
		return errors.New("invalid email format")
	}
	if r.Password == "" {
		return errors.New("password is required")
	}
	if len(r.Password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}
	if r.Password != r.PasswordConfirmation {
		return errors.New("password and password confirmation do not match")
	}
	return nil
}

func isValidEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
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
