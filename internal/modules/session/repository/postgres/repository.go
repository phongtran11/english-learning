package postgres

import (
	"english-learning/internal/modules/session/domain"

	"gorm.io/gorm"
)

type SessionRepository struct {
	db *gorm.DB
}

func NewSessionRepository(db *gorm.DB) domain.SessionRepository {
	return &SessionRepository{db: db}
}

func (r *SessionRepository) Create(session *domain.Session) error {
	sessionModel := FromDomainSession(session)
	if err := r.db.Create(sessionModel).Error; err != nil {
		return err
	}
	session.ID = sessionModel.ID
	session.CreatedAt = sessionModel.CreatedAt
	session.UpdatedAt = sessionModel.UpdatedAt
	return nil
}

func (r *SessionRepository) FindByID(id uint) (*domain.Session, error) {
	var sessionModel SessionGorm
	err := r.db.Preload("User").First(&sessionModel, id).Error
	if err != nil {
		return nil, err
	}
	return sessionModel.ToDomain(), nil
}

func (r *SessionRepository) FindByRefreshToken(refreshToken string) (*domain.Session, error) {
	var sessionModel SessionGorm
	err := r.db.Preload("User").Where("refresh_token = ?", refreshToken).First(&sessionModel).Error
	if err != nil {
		return nil, err
	}
	return sessionModel.ToDomain(), nil
}

func (r *SessionRepository) Revoke(id uint) error {
	return r.db.Model(&SessionGorm{}).Where("id = ?", id).Update("is_revoked", true).Error
}

func (r *SessionRepository) RevokeAllForUser(userID uint) error {
	return r.db.Model(&SessionGorm{}).Where("user_id = ?", userID).Update("is_revoked", true).Error
}

func (r *SessionRepository) Delete(id uint) error {
	return r.db.Delete(&SessionGorm{}, id).Error
}
