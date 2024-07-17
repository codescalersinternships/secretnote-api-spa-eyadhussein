package storage

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/codescalersinternships/secretnote-api-spa-eyadhussein/pkg/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MySQL struct {
	db *gorm.DB
}

func NewMySQL() *MySQL {
	return &MySQL{}
}

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

func (m *MySQL) connect() error {
	err := m.ensureDatabaseExists()
	if err != nil {
		return err
	}
	db, err := gorm.Open(mysql.Open(generateDatasourceName()), &gorm.Config{})
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
	dsn := generateDatasourceNameWithoutDB()

	sqlDB, err := sql.Open("mysql", dsn)
	if err != nil {
		return err
	}
	defer sqlDB.Close()

	_, err = sqlDB.Exec(fmt.Sprintf(sqlStatement, os.Getenv("DB_NAME")))
	if err != nil {
		return err
	}

	log.Printf("database '%s' is ensured to exist or created.", os.Getenv("DB_NAME"))
	return nil
}

func generateDatasourceName() string {
	username := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", username, password, host, port, dbName)
}

func generateDatasourceNameWithoutDB() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/",
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"))
}

func (m *MySQL) migrate() error {
	err := m.db.AutoMigrate(&models.User{}, &models.Note{})
	if err != nil {
		return err
	}
	return nil
}
