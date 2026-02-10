package domain

type SessionRepository interface {
	Create(session *Session) error
	FindByID(id uint) (*Session, error)
	FindByRefreshToken(refreshToken string) (*Session, error)
	Revoke(id uint) error
	RevokeAllForUser(userID uint) error
	Delete(id uint) error
}
