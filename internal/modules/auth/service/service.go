package service

import (
	authDomain "english-learning/internal/modules/auth/domain"
	sessionDomain "english-learning/internal/modules/session/domain"
	userDomain "english-learning/internal/modules/user/domain"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// Service implements authDomain.AuthService.
type Service struct {
	userRepo    userDomain.UserRepository
	sessionRepo sessionDomain.SessionRepository
	jwtSecret   string
}

// NewService creates a new auth Service.
func NewService(userRepo userDomain.UserRepository, sessionRepo sessionDomain.SessionRepository, jwtSecret string) *Service {
	return &Service{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
		jwtSecret:   jwtSecret,
	}
}

func (s *Service) Register(req *authDomain.RegisterRequest) error {
	existingUser, err := s.userRepo.FindByEmail(req.Email)
	if err == nil && existingUser != nil {
		return errors.New("email already registered")
	}

	if err != nil && !errors.Is(err, userDomain.ErrUserNotFound) {
		return fmt.Errorf("checking existing user: %w", err)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hashing password: %w", err)
	}

	user := &userDomain.User{
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	if err := s.userRepo.Create(user); err != nil {
		return fmt.Errorf("creating user: %w", err)
	}

	return nil
}

func (s *Service) Login(req *authDomain.LoginRequest, ip, userAgent string) (*authDomain.TokenPair, error) {
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	accessToken, err := s.generateAccessToken(user)
	if err != nil {
		return nil, fmt.Errorf("generating access token: %w", err)
	}

	refreshToken, err := s.generateRefreshToken(user)
	if err != nil {
		return nil, fmt.Errorf("generating refresh token: %w", err)
	}

	session := &sessionDomain.Session{
		UserID:       user.ID,
		RefreshToken: refreshToken,
		UserAgent:    userAgent,
		ClientIP:     ip,
		ExpiresAt:    time.Now().Add(7 * 24 * time.Hour),
	}

	if err := s.sessionRepo.Create(session); err != nil {
		return nil, fmt.Errorf("creating session: %w", err)
	}

	return &authDomain.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *Service) RefreshToken(refreshToken string) (*authDomain.TokenPair, error) {
	// Verify refresh token
	token, err := jwt.ParseWithClaims(refreshToken, &authClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.jwtSecret), nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid refresh token")
	}

	// Check if session exists and is valid
	session, err := s.sessionRepo.FindByRefreshToken(refreshToken)
	if err != nil || session == nil {
		return nil, errors.New("invalid session")
	}

	if session.IsRevoked {
		return nil, errors.New("session revoked")
	}

	// Revoke current session (Token Rotation)
	if err := s.sessionRepo.Revoke(session.ID); err != nil {
		return nil, fmt.Errorf("revoking session: %w", err)
	}

	// Check if associated user exists
	user, err := s.userRepo.FindByID(session.UserID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Generate new access token and new refresh token
	accessToken, err := s.generateAccessToken(user)
	if err != nil {
		return nil, fmt.Errorf("generating access token: %w", err)
	}

	newRefreshToken, err := s.generateRefreshToken(user)
	if err != nil {
		return nil, fmt.Errorf("generating refresh token: %w", err)
	}

	// Create new session with the new refresh token
	newSession := &sessionDomain.Session{
		UserID:       user.ID,
		RefreshToken: newRefreshToken,
		UserAgent:    session.UserAgent,
		ClientIP:     session.ClientIP,
		ExpiresAt:    time.Now().Add(7 * 24 * time.Hour),
	}

	if err := s.sessionRepo.Create(newSession); err != nil {
		return nil, fmt.Errorf("creating new session: %w", err)
	}

	return &authDomain.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

func (s *Service) generateAccessToken(user *userDomain.User) (string, error) {
	claims := authClaims{
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   strconv.FormatUint(uint64(user.ID), 10),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			Issuer:    "english-learning",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}

func (s *Service) generateRefreshToken(user *userDomain.User) (string, error) {
	claims := authClaims{
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   strconv.FormatUint(uint64(user.ID), 10),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
			Issuer:    "english-learning",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}

func (s *Service) Logout(refreshToken string) error {
	session, err := s.sessionRepo.FindByRefreshToken(refreshToken)
	if err != nil || session == nil {
		return nil // Already logged out or invalid
	}

	if err := s.sessionRepo.Revoke(session.ID); err != nil {
		return fmt.Errorf("revoking session: %w", err)
	}

	return nil
}

func (s *Service) LogoutAll(userID uint) error {
	if err := s.sessionRepo.RevokeAllForUser(userID); err != nil {
		return fmt.Errorf("revoking all sessions: %w", err)
	}

	return nil
}
