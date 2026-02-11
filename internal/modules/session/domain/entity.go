package domain

import (
	"time"
)

type Session struct {
	ID           uint
	UserID       uint
	RefreshToken string
	UserAgent    string
	ClientIP     string
	IsRevoked    bool
	ExpiresAt    time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
