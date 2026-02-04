package userservice

import (
	"english-learning/internal/core/domain"
	"english-learning/internal/core/ports"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo ports.UserRepository
}

func NewService(repo ports.UserRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(req *domain.RegisterRequest) error {
	existing, _ := s.repo.FindByEmail(req.Email)
	if existing != nil {
		return errors.New("email already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &domain.User{
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	return s.repo.Create(user)
}

func (s *Service) Get(id uint) (*domain.User, error) {
	return s.repo.FindByID(id)
}

func (s *Service) Update(user *domain.User) error {
	if user.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		user.Password = string(hashedPassword)
	}
	return s.repo.Update(user)
}

func (s *Service) Delete(id uint) error {
	return s.repo.Delete(id)
}

func (s *Service) List(page, pageSize int) ([]domain.User, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize
	return s.repo.List(offset, pageSize)
}
