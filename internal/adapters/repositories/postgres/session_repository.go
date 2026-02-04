package postgres

import (
	"english-learning/internal/core/domain"
	"english-learning/internal/core/ports"

	"gorm.io/gorm"
)

type SessionRepository struct {
	db *gorm.DB
}

func NewSessionRepository(db *gorm.DB) ports.SessionRepository {
	return &SessionRepository{db: db}
}

func (r *SessionRepository) Create(session *domain.Session) error {
	return r.db.Create(session).Error
}

func (r *SessionRepository) GetByRefreshToken(refreshToken string) (*domain.Session, error) {
	var session domain.Session
	err := r.db.Where("refresh_token = ?", refreshToken).First(&session).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *SessionRepository) Revoke(id uint) error {
	return r.db.Model(&domain.Session{}).Where("id = ?", id).Update("is_revoked", true).Error
}

func (r *SessionRepository) RevokeByUserID(userID uint) error {
	return r.db.Model(&domain.Session{}).Where("user_id = ?", userID).Update("is_revoked", true).Error
}
