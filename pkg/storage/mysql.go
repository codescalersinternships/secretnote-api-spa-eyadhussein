package storage

import (
	"log"

	"github.com/codescalersinternships/secretnote-api-spa-eyadhussein/pkg/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MySQL struct {
	db *gorm.DB
}

func NewMySQL() (*MySQL, error) {
	m := &MySQL{}

	err := m.init()

	if err != nil {
		return nil, err
	}
	return m, nil
}

func (m *MySQL) init() error {
	err := m.connect()
	if err != nil {
		return err
	}
	err = m.createUserTable()

	if err != nil {
		return err
	}

	return nil
}

func (m *MySQL) connect() error {
	dsn := "root:123@tcp(127.0.0.1:3306)/secretnote?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	m.db = db

	log.Println("successfully connected to the database")
	return nil
}

func (m *MySQL) createUserTable() error {
	err := m.db.AutoMigrate(&models.User{})
	if err != nil {
		return err
	}
	return nil
}

func (m *MySQL) CreateUser(user *models.User) error {
	result := m.db.Create(user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (m *MySQL) GetUserByUsername(username string) (*models.User, error) {
	user := &models.User{}
	result := m.db.Where("username = ?", username).First(user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}
