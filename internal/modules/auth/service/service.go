package service

import (
	authDomain "english-learning/internal/modules/auth/domain"
	sessionDomain "english-learning/internal/modules/session/domain"
	userDomain "english-learning/internal/modules/user/domain"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	userRepo    userDomain.UserRepository
	sessionRepo sessionDomain.SessionRepository
	jwtSecret   string
}

func NewService(userRepo userDomain.UserRepository, sessionRepo sessionDomain.SessionRepository, jwtSecret string) *Service {

	return &Service{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
		jwtSecret:   jwtSecret,
	}
}

func (s *Service) Register(req *authDomain.RegisterRequest) error {
	existingUser, _ := s.userRepo.FindByEmail(req.Email)
	if existingUser != nil {
		return errors.New("email already registered")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &userDomain.User{
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	return s.userRepo.Create(user)
}

func (s *Service) Login(req *authDomain.LoginRequest, ip, userAgent string) (*authDomain.TokenPair, error) {
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Create tokens
	accessToken, err := s.generateAccessToken(user)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.generateRefreshToken(user)
	if err != nil {
		return nil, err
	}

	// Create Session
	session := &sessionDomain.Session{
		UserID:       user.ID,
		RefreshToken: refreshToken,
		UserAgent:    userAgent,
		ClientIP:     ip,
		ExpiresAt:    time.Now().Add(7 * 24 * time.Hour), // 7 days
	}

	if err := s.sessionRepo.Create(session); err != nil {
		return nil, err
	}

	return &authDomain.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *Service) RefreshToken(refreshToken string) (*authDomain.TokenPair, error) {
	// Verify refresh token
	token, err := jwt.ParseWithClaims(refreshToken, &authDomain.AuthClaims{}, func(token *jwt.Token) (interface{}, error) {
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

	// Check if associated user exists (optional but safe)
	user, err := s.userRepo.FindByID(session.UserID)
	if err != nil {
		return nil, errors.New("user not found")
	}


	return &authDomain.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *Service) generateAccessToken(user *userDomain.User) (string, error) {
	claims := authDomain.AuthClaims{
		UserId: user.ID,
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			Issuer:    "english-learning",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}

func (s *Service) generateRefreshToken(user *userDomain.User) (string, error) {
	claims := authDomain.AuthClaims{
		UserId: user.ID,
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
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
	return s.sessionRepo.Revoke(session.ID)
}

func (s *Service) LogoutAll(userID uint) error {
	return s.sessionRepo.RevokeAllForUser(userID)
}
