package http

import (
	authDomain "english-learning/internal/modules/auth/domain"

	"github.com/stretchr/testify/mock"
)

// MockAuthService is a mock implementation of authDomain.AuthService.
type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) Register(req *authDomain.RegisterRequest) error {
	args := m.Called(req)
	return args.Error(0)
}

func (m *MockAuthService) Login(req *authDomain.LoginRequest, ip, userAgent string) (*authDomain.TokenPair, error) {
	args := m.Called(req, ip, userAgent)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*authDomain.TokenPair), args.Error(1)
}

func (m *MockAuthService) RefreshToken(refreshToken string) (*authDomain.TokenPair, error) {
	args := m.Called(refreshToken)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*authDomain.TokenPair), args.Error(1)
}

func (m *MockAuthService) Logout(refreshToken string) error {
	args := m.Called(refreshToken)
	return args.Error(0)
}

func (m *MockAuthService) LogoutAll(userID uint) error {
	args := m.Called(userID)
	return args.Error(0)
}
