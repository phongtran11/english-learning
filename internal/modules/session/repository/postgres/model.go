package postgres

import (
	"english-learning/internal/modules/session/domain"
	userPostgres "english-learning/internal/modules/user/repository/postgres"
	"time"
)

type Session struct {
	ID           uint      `gorm:"primaryKey"`
	UserID       uint      `gorm:"not null;index"`
	RefreshToken string    `gorm:"type:text;not null"`
	UserAgent    string    `gorm:"type:text"`
	ClientIP     string    `gorm:"type:varchar(45)"`
	IsRevoked    bool      `gorm:"not null;default:false"`
	ExpiresAt    time.Time `gorm:"not null"`
	CreatedAt    time.Time
	UpdatedAt    time.Time

	// Relation (Belongs To)
	User *userPostgres.User `gorm:"constraint:OnDelete:CASCADE;"`
}

func (m *Session) ToDomain() *domain.Session {
    if m == nil {
        return nil
    }
	return &domain.Session{
		ID:           m.ID,
		UserID:       m.UserID,
		RefreshToken: m.RefreshToken,
		UserAgent:    m.UserAgent,
		ClientIP:     m.ClientIP,
		IsRevoked:    m.IsRevoked,
		ExpiresAt:    m.ExpiresAt,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
	}
}

func FromDomainSession(s *domain.Session) *Session {
    if s == nil {
        return nil
    }
	return &Session{
		ID:           s.ID,
		UserID:       s.UserID,
		RefreshToken: s.RefreshToken,
		UserAgent:    s.UserAgent,
		ClientIP:     s.ClientIP,
		IsRevoked:    s.IsRevoked,
		ExpiresAt:    s.ExpiresAt,
		CreatedAt:    s.CreatedAt,
		UpdatedAt:    s.UpdatedAt,
	}
}
