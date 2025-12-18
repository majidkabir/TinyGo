package service

import (
	"fmt"

	"github.com/majidkabir/TinyGo/internal/models"
	"github.com/majidkabir/TinyGo/internal/repository"
)

// UserService handles business logic for users
type UserService struct {
	repo *repository.UserRepository
}

// NewUserService creates a new user service
func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

// CreateUser creates a new user with validation
func (s *UserService) CreateUser(req *models.CreateUserRequest) (*models.User, error) {
	// Business rule: Check if email already exists
	existingUser, err := s.repo.GetByEmail(req.Email)
	if err == nil && existingUser != nil {
		return nil, fmt.Errorf("email already exists")
	}

	// Business rule: Additional validation
	if req.Age < 18 {
		return nil, fmt.Errorf("user must be at least 18 years old")
	}

	// Create user
	user, err := s.repo.Create(req)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

// GetUser retrieves a user by ID
func (s *UserService) GetUser(id int) (*models.User, error) {
	user, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// ListUsers retrieves all users with pagination
func (s *UserService) ListUsers(page, pageSize int) ([]*models.User, int, error) {
	// Calculate offset
	offset := (page - 1) * pageSize

	// Get users
	users, err := s.repo.GetAll(pageSize, offset)
	if err != nil {
		return nil, 0, err
	}

	// Get total count for pagination
	total, err := s.repo.Count()
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// UpdateUser updates a user
func (s *UserService) UpdateUser(id int, req *models.UpdateUserRequest) (*models.User, error) {
	// Check if user exists
	_, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Business rule: If updating email, check if new email exists
	if req.Email != "" {
		existingUser, err := s.repo.GetByEmail(req.Email)
		if err == nil && existingUser != nil && existingUser.ID != id {
			return nil, fmt.Errorf("email already exists")
		}
	}

	// Update user
	user, err := s.repo.Update(id, req)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return user, nil
}

// DeleteUser deletes a user
func (s *UserService) DeleteUser(id int) error {
	// Check if user exists
	_, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	// Delete user
	if err := s.repo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}
