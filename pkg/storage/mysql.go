package storage

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/codescalersinternships/secretnote-api-spa-eyadhussein/pkg/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Config is a configuration for the MySQL storage
type Config struct {
	username string
	password string
	host     string
	port     string
	dbName   string
}

// NewConfig creates a new MySQL storage configuration
func NewConfig(username, password, host, port, dbName string) *Config {
	return &Config{
		username: username,
		password: password,
		host:     host,
		port:     port,
		dbName:   dbName,
	}
}

func (c *Config) generateDatasource(includeName bool) string {
	if includeName {
		return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", c.username, c.password, c.host, c.port, c.dbName)
	} else {
		return fmt.Sprintf("%s:%s@tcp(%s:%s)/", c.username, c.password, c.host, c.port)
	}
}

// MySQL is a storage implementation that uses MySQL
type MySQL struct {
	db     *gorm.DB
	config *Config
}

// NewMySQL creates a new MySQL storage
func NewMySQL(config *Config) *MySQL {
	return &MySQL{
		config: config,
	}
}

// Init initializes the MySQL storage
func (m *MySQL) Init() error {
	err := m.connect()
	if err != nil {
		return err
	}
	err = m.migrate()

	if err != nil {
		return err
	}

	return nil
}

// CreateUser creates a new user
func (m *MySQL) CreateUser(user *models.User) error {
	tx := m.db.Begin()
	if err := tx.Create(user).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to create user %w", err)
	}

	tx.Commit()
	return nil
}

// GetUserByUsername gets a user by username
func (m *MySQL) GetUserByUsername(username string) (*models.User, error) {
	var user *models.User
	result := m.db.Where("username = ?", username).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}

// CreateNote creates a new note
func (m *MySQL) CreateNote(note *models.Note) error {
	tx := m.db.Begin()
	tx.Create(note)
	if tx.Error != nil {
		tx.Rollback()
		return tx.Error
	}

	tx.Commit()
	return nil
}

// GetNoteByID gets a note by ID
func (m *MySQL) GetNoteByID(id string) (*models.Note, error) {
	var note *models.Note
	result := m.db.Where("id = ?", id).First(&note)
	if result.Error != nil {
		return nil, result.Error
	}
	return note, nil
}

// GetNotesByUserID gets notes by user ID
func (m *MySQL) GetNotesByUserID(userID uint) ([]*models.Note, error) {
	var notes []*models.Note
	result := m.db.Where("user_id = ?", userID).Find(&notes)
	if result.Error != nil {
		return nil, result.Error
	}
	return notes, nil
}

// IncrementNoteViews increments the views of a note
func (m *MySQL) IncrementNoteViews(id string) error {
	tx := m.db.Begin()
	var note *models.Note
	result := tx.Where("id = ?", id).First(&note)
	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	note.CurrentViews++
	result = tx.Save(note)
	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}
	tx.Commit()
	return nil
}

// DeleteNoteByID deletes a note by ID
func (m *MySQL) DeleteNoteByID(id string) error {
	tx := m.db.Begin()
	log.Println("Deleting note with ID: ", id)
	result := tx.Where("id = ?", id).Delete(&models.Note{})
	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}
	tx.Commit()
	return nil
}

// Clear clears the database
func (m *MySQL) Clear() error {
	tx := m.db.Begin()
	tx.Exec("DROP DATABASE " + m.config.dbName)
	if tx.Error != nil {
		tx.Rollback()
		return tx.Error
	}
	tx.Exec("CREATE DATABASE " + m.config.dbName)
	if tx.Error != nil {
		tx.Rollback()
		return tx.Error
	}
	tx.Commit()
	return nil
}

func (m *MySQL) connect() error {
	err := m.ensureDatabaseExists()
	if err != nil {
		return err
	}
	db, err := gorm.Open(mysql.Open(m.config.generateDatasource(true)), &gorm.Config{})
	if err != nil {
		return err
	}

	m.db = db

	log.Println("successfully connected to the database")
	return nil
}

func (m *MySQL) ensureDatabaseExists() error {
	var sqlStatement = `
		CREATE DATABASE IF NOT EXISTS %s;
	`
	dsn := m.config.generateDatasource(false)

	sqlDB, err := sql.Open("mysql", dsn)
	if err != nil {
		return err
	}
	defer sqlDB.Close()

	_, err = sqlDB.Exec(fmt.Sprintf(sqlStatement, m.config.dbName))
	if err != nil {
		return err
	}
	return nil
}

func (m *MySQL) migrate() error {
	err := m.db.AutoMigrate(&models.User{}, &models.Note{})
	if err != nil {
		return err
	}
	return nil
}
