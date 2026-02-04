package ports

import "english-learning/internal/core/domain"

type SessionRepository interface {
	Create(session *domain.Session) error
	GetByRefreshToken(refreshToken string) (*domain.Session, error)
	Revoke(id uint) error
	RevokeByUserID(userID uint) error
}
