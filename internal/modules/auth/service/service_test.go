package service

import (
	authDomain "english-learning/internal/modules/auth/domain"
	sessionDomain "english-learning/internal/modules/session/domain"
	userDomain "english-learning/internal/modules/user/domain"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

const testJWTSecret = "test-secret-key-for-unit-tests"

// newTestService creates a Service with mock dependencies for testing.
func newTestService() (*Service, *MockUserRepository, *MockSessionRepository) {
	userRepo := new(MockUserRepository)
	sessionRepo := new(MockSessionRepository)
	svc := NewService(userRepo, sessionRepo, testJWTSecret)
	return svc, userRepo, sessionRepo
}

// hashPassword is a test helper to create a bcrypt hashed password.
func hashPassword(t *testing.T, password string) string {
	t.Helper()
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}
	return string(hashed)
}

// --- Register Tests ---

func TestRegister_Success(t *testing.T) {
	t.Parallel()
	svc, userRepo, _ := newTestService()

	req := &authDomain.RegisterRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	userRepo.On("FindByEmail", req.Email).Return(nil, userDomain.ErrUserNotFound)
	userRepo.On("Create", mock.AnythingOfType("*domain.User")).Return(nil)

	err := svc.Register(req)

	assert.NoError(t, err)
	userRepo.AssertExpectations(t)
}

func TestRegister_EmailAlreadyExists(t *testing.T) {
	t.Parallel()
	svc, userRepo, _ := newTestService()

	existingUser := &userDomain.User{
		ID:    1,
		Email: "test@example.com",
	}

	req := &authDomain.RegisterRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	userRepo.On("FindByEmail", req.Email).Return(existingUser, nil)

	err := svc.Register(req)

	assert.Error(t, err)
	assert.Equal(t, "email already registered", err.Error())
	userRepo.AssertNotCalled(t, "Create", mock.Anything)
}

func TestRegister_RepositoryError(t *testing.T) {
	t.Parallel()
	svc, userRepo, _ := newTestService()

	req := &authDomain.RegisterRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	dbErr := errors.New("database connection failed")
	userRepo.On("FindByEmail", req.Email).Return(nil, dbErr)

	err := svc.Register(req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "checking existing user")
	userRepo.AssertNotCalled(t, "Create", mock.Anything)
}

func TestRegister_CreateUserError(t *testing.T) {
	t.Parallel()
	svc, userRepo, _ := newTestService()

	req := &authDomain.RegisterRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	userRepo.On("FindByEmail", req.Email).Return(nil, userDomain.ErrUserNotFound)
	userRepo.On("Create", mock.AnythingOfType("*domain.User")).Return(errors.New("insert failed"))

	err := svc.Register(req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "creating user")
}

// --- Login Tests ---

func TestLogin_Success(t *testing.T) {
	t.Parallel()
	svc, userRepo, sessionRepo := newTestService()

	password := "password123"
	user := &userDomain.User{
		ID:       1,
		Email:    "test@example.com",
		Password: hashPassword(t, password),
	}

	req := &authDomain.LoginRequest{
		Email:    "test@example.com",
		Password: password,
	}

	userRepo.On("FindByEmail", req.Email).Return(user, nil)
	sessionRepo.On("Create", mock.AnythingOfType("*domain.Session")).Return(nil)

	tokenPair, err := svc.Login(req, "127.0.0.1", "TestAgent/1.0")

	assert.NoError(t, err)
	assert.NotNil(t, tokenPair)
	assert.NotEmpty(t, tokenPair.AccessToken)
	assert.NotEmpty(t, tokenPair.RefreshToken)
	userRepo.AssertExpectations(t)
	sessionRepo.AssertExpectations(t)
}

func TestLogin_UserNotFound(t *testing.T) {
	t.Parallel()
	svc, userRepo, _ := newTestService()

	req := &authDomain.LoginRequest{
		Email:    "nonexistent@example.com",
		Password: "password123",
	}

	userRepo.On("FindByEmail", req.Email).Return(nil, userDomain.ErrUserNotFound)

	tokenPair, err := svc.Login(req, "127.0.0.1", "TestAgent/1.0")

	assert.Error(t, err)
	assert.Nil(t, tokenPair)
	assert.Equal(t, "invalid credentials", err.Error())
}

func TestLogin_InvalidPassword(t *testing.T) {
	t.Parallel()
	svc, userRepo, _ := newTestService()

	user := &userDomain.User{
		ID:       1,
		Email:    "test@example.com",
		Password: hashPassword(t, "correct-password"),
	}

	req := &authDomain.LoginRequest{
		Email:    "test@example.com",
		Password: "wrong-password",
	}

	userRepo.On("FindByEmail", req.Email).Return(user, nil)

	tokenPair, err := svc.Login(req, "127.0.0.1", "TestAgent/1.0")

	assert.Error(t, err)
	assert.Nil(t, tokenPair)
	assert.Equal(t, "invalid credentials", err.Error())
	// Session should NOT be created
	sessionRepo := new(MockSessionRepository)
	sessionRepo.AssertNotCalled(t, "Create", mock.Anything)
}

func TestLogin_SessionCreateError(t *testing.T) {
	t.Parallel()
	svc, userRepo, sessionRepo := newTestService()

	password := "password123"
	user := &userDomain.User{
		ID:       1,
		Email:    "test@example.com",
		Password: hashPassword(t, password),
	}

	req := &authDomain.LoginRequest{
		Email:    "test@example.com",
		Password: password,
	}

	userRepo.On("FindByEmail", req.Email).Return(user, nil)
	sessionRepo.On("Create", mock.AnythingOfType("*domain.Session")).Return(errors.New("session insert failed"))

	tokenPair, err := svc.Login(req, "127.0.0.1", "TestAgent/1.0")

	assert.Error(t, err)
	assert.Nil(t, tokenPair)
	assert.Contains(t, err.Error(), "creating session")
}

// --- RefreshToken Tests ---

func TestRefreshToken_Success(t *testing.T) {
	t.Parallel()
	svc, userRepo, sessionRepo := newTestService()

	// Generate a valid refresh token first
	user := &userDomain.User{ID: 1, Email: "test@example.com"}
	validRefreshToken, err := svc.generateRefreshToken(user)
	assert.NoError(t, err)

	session := &sessionDomain.Session{
		ID:           1,
		UserID:       1,
		RefreshToken: validRefreshToken,
		UserAgent:    "TestAgent/1.0",
		ClientIP:     "127.0.0.1",
		IsRevoked:    false,
		ExpiresAt:    time.Now().Add(7 * 24 * time.Hour),
	}

	sessionRepo.On("FindByRefreshToken", validRefreshToken).Return(session, nil)
	sessionRepo.On("Revoke", uint(1)).Return(nil)
	userRepo.On("FindByID", uint(1)).Return(user, nil)
	sessionRepo.On("Create", mock.AnythingOfType("*domain.Session")).Return(nil)

	tokenPair, err := svc.RefreshToken(validRefreshToken)

	assert.NoError(t, err)
	assert.NotNil(t, tokenPair)
	assert.NotEmpty(t, tokenPair.AccessToken)
	assert.NotEmpty(t, tokenPair.RefreshToken)
	sessionRepo.AssertExpectations(t)
	userRepo.AssertExpectations(t)
}

func TestRefreshToken_InvalidToken(t *testing.T) {
	t.Parallel()
	svc, _, _ := newTestService()

	tokenPair, err := svc.RefreshToken("invalid-token-string")

	assert.Error(t, err)
	assert.Nil(t, tokenPair)
	assert.Equal(t, "invalid refresh token", err.Error())
}

func TestRefreshToken_SessionNotFound(t *testing.T) {
	t.Parallel()
	svc, _, sessionRepo := newTestService()

	user := &userDomain.User{ID: 1, Email: "test@example.com"}
	validRefreshToken, _ := svc.generateRefreshToken(user)

	sessionRepo.On("FindByRefreshToken", validRefreshToken).Return(nil, errors.New("not found"))

	tokenPair, err := svc.RefreshToken(validRefreshToken)

	assert.Error(t, err)
	assert.Nil(t, tokenPair)
	assert.Equal(t, "invalid session", err.Error())
}

func TestRefreshToken_SessionRevoked(t *testing.T) {
	t.Parallel()
	svc, _, sessionRepo := newTestService()

	user := &userDomain.User{ID: 1, Email: "test@example.com"}
	validRefreshToken, _ := svc.generateRefreshToken(user)

	session := &sessionDomain.Session{
		ID:           1,
		UserID:       1,
		RefreshToken: validRefreshToken,
		IsRevoked:    true,
	}

	sessionRepo.On("FindByRefreshToken", validRefreshToken).Return(session, nil)

	tokenPair, err := svc.RefreshToken(validRefreshToken)

	assert.Error(t, err)
	assert.Nil(t, tokenPair)
	assert.Equal(t, "session revoked", err.Error())
	// Should NOT call Revoke again
	sessionRepo.AssertNotCalled(t, "Revoke", mock.Anything)
}

func TestRefreshToken_UserNotFound(t *testing.T) {
	t.Parallel()
	svc, userRepo, sessionRepo := newTestService()

	user := &userDomain.User{ID: 1, Email: "test@example.com"}
	validRefreshToken, _ := svc.generateRefreshToken(user)

	session := &sessionDomain.Session{
		ID:           1,
		UserID:       1,
		RefreshToken: validRefreshToken,
		IsRevoked:    false,
	}

	sessionRepo.On("FindByRefreshToken", validRefreshToken).Return(session, nil)
	sessionRepo.On("Revoke", uint(1)).Return(nil)
	userRepo.On("FindByID", uint(1)).Return(nil, errors.New("user not found"))

	tokenPair, err := svc.RefreshToken(validRefreshToken)

	assert.Error(t, err)
	assert.Nil(t, tokenPair)
	assert.Equal(t, "user not found", err.Error())
}

// --- Logout Tests ---

func TestLogout_Success(t *testing.T) {
	t.Parallel()
	svc, _, sessionRepo := newTestService()

	user := &userDomain.User{ID: 1, Email: "test@example.com"}
	validRefreshToken, _ := svc.generateRefreshToken(user)

	session := &sessionDomain.Session{
		ID:           1,
		UserID:       1,
		RefreshToken: validRefreshToken,
	}

	sessionRepo.On("FindByRefreshToken", validRefreshToken).Return(session, nil)
	sessionRepo.On("Revoke", uint(1)).Return(nil)

	err := svc.Logout(validRefreshToken)

	assert.NoError(t, err)
	sessionRepo.AssertExpectations(t)
}

func TestLogout_SessionNotFound_Idempotent(t *testing.T) {
	t.Parallel()
	svc, _, sessionRepo := newTestService()

	sessionRepo.On("FindByRefreshToken", "some-token").Return(nil, errors.New("not found"))

	err := svc.Logout("some-token")

	// Should not return error — idempotent behavior
	assert.NoError(t, err)
}

func TestLogout_RevokeError(t *testing.T) {
	t.Parallel()
	svc, _, sessionRepo := newTestService()

	user := &userDomain.User{ID: 1, Email: "test@example.com"}
	validRefreshToken, _ := svc.generateRefreshToken(user)

	session := &sessionDomain.Session{
		ID:           1,
		UserID:       1,
		RefreshToken: validRefreshToken,
	}

	sessionRepo.On("FindByRefreshToken", validRefreshToken).Return(session, nil)
	sessionRepo.On("Revoke", uint(1)).Return(errors.New("revoke failed"))

	err := svc.Logout(validRefreshToken)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "revoking session")
}

// --- LogoutAll Tests ---

func TestLogoutAll_Success(t *testing.T) {
	t.Parallel()
	svc, _, sessionRepo := newTestService()

	sessionRepo.On("RevokeAllForUser", uint(1)).Return(nil)

	err := svc.LogoutAll(1)

	assert.NoError(t, err)
	sessionRepo.AssertExpectations(t)
}

func TestLogoutAll_Error(t *testing.T) {
	t.Parallel()
	svc, _, sessionRepo := newTestService()

	sessionRepo.On("RevokeAllForUser", uint(1)).Return(errors.New("db error"))

	err := svc.LogoutAll(1)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "revoking all sessions")
}
