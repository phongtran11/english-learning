package service

import (
	"english-learning/internal/modules/user/domain"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// Service implements domain.UserService.
type Service struct {
	repo domain.UserRepository
}

// NewService creates a new user Service.
func NewService(repo domain.UserRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(req *domain.User) error {
	existing, err := s.repo.FindByEmail(req.Email)
	if err == nil && existing != nil {
		return errors.New("email already exists")
	}

	if err != nil && !errors.Is(err, domain.ErrUserNotFound) {
		return fmt.Errorf("checking existing user: %w", err)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hashing password: %w", err)
	}

	req.Password = string(hashedPassword)

	if err := s.repo.Create(req); err != nil {
		return fmt.Errorf("creating user: %w", err)
	}

	return nil
}

func (s *Service) Get(id uint) (*domain.User, error) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("finding user by id: %w", err)
	}

	return user, nil
}

func (s *Service) Update(user *domain.User) error {
	if user.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("hashing password: %w", err)
		}
		user.Password = string(hashedPassword)
	}

	if err := s.repo.Update(user); err != nil {
		return fmt.Errorf("updating user: %w", err)
	}

	return nil
}

func (s *Service) Delete(id uint) error {
	if err := s.repo.Delete(id); err != nil {
		return fmt.Errorf("deleting user: %w", err)
	}

	return nil
}

func (s *Service) List(page, pageSize int) ([]domain.User, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	users, count, err := s.repo.List(offset, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("listing users: %w", err)
	}

	return users, count, nil
}
