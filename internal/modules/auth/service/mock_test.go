package service

import (
	sessionDomain "english-learning/internal/modules/session/domain"
	userDomain "english-learning/internal/modules/user/domain"

	"github.com/stretchr/testify/mock"
)

// MockUserRepository is a mock implementation of userDomain.UserRepository.
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(user *userDomain.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) FindByEmail(email string) (*userDomain.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*userDomain.User), args.Error(1)
}

func (m *MockUserRepository) FindByID(id uint) (*userDomain.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*userDomain.User), args.Error(1)
}

func (m *MockUserRepository) Update(user *userDomain.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockUserRepository) List(offset, limit int) ([]userDomain.User, int64, error) {
	args := m.Called(offset, limit)
	return args.Get(0).([]userDomain.User), args.Get(1).(int64), args.Error(2)
}

// MockSessionRepository is a mock implementation of sessionDomain.SessionRepository.
type MockSessionRepository struct {
	mock.Mock
}

func (m *MockSessionRepository) Create(session *sessionDomain.Session) error {
	args := m.Called(session)
	return args.Error(0)
}

func (m *MockSessionRepository) FindByID(id uint) (*sessionDomain.Session, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*sessionDomain.Session), args.Error(1)
}

func (m *MockSessionRepository) FindByRefreshToken(refreshToken string) (*sessionDomain.Session, error) {
	args := m.Called(refreshToken)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*sessionDomain.Session), args.Error(1)
}

func (m *MockSessionRepository) Revoke(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockSessionRepository) RevokeAllForUser(userID uint) error {
	args := m.Called(userID)
	return args.Error(0)
}

func (m *MockSessionRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}
