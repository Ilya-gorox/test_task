package service

import (
	"database/sql"
	"fmt"
	"time"

	"test_service/internal/model"
	"test_service/internal/repository"

	"github.com/google/uuid"
)

type Service struct {
	repo *repository.Repository
}

func NewService(repo *repository.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateUser(input *model.CreateUserInput) (*model.User, error) {
	user := &model.User{
		ID:        uuid.New(),
		Firstname: input.Firstname,
		Lastname:  input.Lastname,
		Email:     input.Email,
		Age:       input.Age,
		Created:   time.Now(),
	}

	if err := s.repo.CreateUser(user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

func (s *Service) GetUser(id string) (*model.User, error) {
	user, err := s.repo.GetUser(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user with id %s not found", id)
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}

func (s *Service) UpdateUser(id string, input *model.UpdateUserInput) error {
	if err := s.repo.UpdateUser(id, input); err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("user with id %s not found", id)
		}
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}
