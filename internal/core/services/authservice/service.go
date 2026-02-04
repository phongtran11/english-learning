package authservice

import (
	"english-learning/configs"
	"english-learning/internal/core/domain"
	"english-learning/internal/core/ports"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo        ports.UserRepository
	sessionRepo ports.SessionRepository
	jwtConfig   configs.JWTConfig
}

func NewAuthService(repo ports.UserRepository, sessionRepo ports.SessionRepository, jwtConfig configs.JWTConfig) *AuthService {

	return &AuthService{
		repo:        repo,
		sessionRepo: sessionRepo,
		jwtConfig:   jwtConfig,
	}
}

func (s *AuthService) Register(req *domain.RegisterRequest) error {
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

func (s *AuthService) Login(req *domain.LoginRequest, userAgent string, clientIP string) (*domain.TokenPair, error) {
	user, err := s.repo.FindByEmail(req.Email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	tokenPair, err := s.generateTokenPair(user)
	if err != nil {
		return nil, err
	}

	// Create Session
	session := &domain.Session{
		UserID:       user.ID,
		RefreshToken: tokenPair.RefreshToken,
		UserAgent:    userAgent,
		ClientIP:     clientIP,
		ExpiresAt:    time.Now().Add(time.Hour * time.Duration(s.jwtConfig.RefreshExpiryHour)),
	}

	if err := s.sessionRepo.Create(session); err != nil {
		return nil, err
	}

	return tokenPair, nil
}

func (s *AuthService) RefreshToken(refreshTokenStr string) (*domain.TokenPair, error) {
	token, err := jwt.Parse(refreshTokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.jwtConfig.Secret), nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid refresh token")
	}

	// Check Session
	session, err := s.sessionRepo.GetByRefreshToken(refreshTokenStr)
	if err != nil {
		return nil, errors.New("session not found")
	}
	if session.IsRevoked {
		return nil, errors.New("session revoked")
	}
	if time.Now().After(session.ExpiresAt) {
		return nil, errors.New("session expired")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return nil, errors.New("invalid user claim")
	}
	userID := uint(userIDFloat)

	user, err := s.repo.FindByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Optional: Revoke old session and create new one (Rotation)
	// For now, simple standard refresh (keep usage of same refresh token valid until expiry?)
	// If standard flow: usually refresh token is reused until expiry OR rotated.
	// Let's implement rotation: Revoke old, create new.
	if err := s.sessionRepo.Revoke(session.ID); err != nil {
		return nil, err
	}

	newTokenPair, err := s.generateTokenPair(user)
	if err != nil {
		return nil, err
	}

	// Create new session
	newSession := &domain.Session{
		UserID:       user.ID,
		RefreshToken: newTokenPair.RefreshToken,
		UserAgent:    session.UserAgent, // Preserve UA/IP
		ClientIP:     session.ClientIP,
		ExpiresAt:    time.Now().Add(time.Hour * time.Duration(s.jwtConfig.RefreshExpiryHour)),
	}

	if err := s.sessionRepo.Create(newSession); err != nil {
		return nil, err
	}

	return newTokenPair, nil
}

func (s *AuthService) Logout(refreshTokenStr string) error {
	session, err := s.sessionRepo.GetByRefreshToken(refreshTokenStr)
	if err != nil {
		return nil // Already invalid/not found, consider it logout success
	}
	return s.sessionRepo.Revoke(session.ID)
}

func (s *AuthService) generateTokenPair(user *domain.User) (*domain.TokenPair, error) {
	// Access Token
	atClaims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     time.Now().Add(time.Hour * time.Duration(s.jwtConfig.AccessExpiryHour)).Unix(),
	}
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	accessToken, err := at.SignedString([]byte(s.jwtConfig.Secret))
	if err != nil {
		return nil, err
	}

	// Refresh Token
	rtClaims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * time.Duration(s.jwtConfig.RefreshExpiryHour)).Unix(),
	}
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	refreshToken, err := rt.SignedString([]byte(s.jwtConfig.Secret))
	if err != nil {
		return nil, err
	}

	return &domain.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
