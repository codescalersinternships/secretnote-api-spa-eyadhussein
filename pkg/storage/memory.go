package storage

import "github.com/codescalersinternships/secretnote-api-spa-eyadhussein/pkg/models"

type Memory struct{}

func NewMemory() *Memory {
	return &Memory{}
}

func (m *Memory) GetUserByID(id int) (*models.User, error) {
	return &models.User{
		ID:       1,
		Username: "test",
		Password: "test",
		Email:    "test@gmail.com",
	}, nil
}

func (m *Memory) GetUserByUsername(username string) (*models.User, error) {
	return nil, nil
}
