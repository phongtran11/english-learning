package domain

import (
	"time"
)

type Session struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	UserID       uint      `gorm:"not null;index" json:"user_id"`
	RefreshToken string    `gorm:"type:text;not null" json:"refresh_token"`
	UserAgent    string    `gorm:"type:text" json:"user_agent"`
	ClientIP     string    `gorm:"type:varchar(45)" json:"client_ip"`
	IsRevoked    bool      `gorm:"not null;default:false" json:"is_revoked"`
	ExpiresAt    time.Time `gorm:"not null" json:"expires_at"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	// Relation (Belongs To)
	User *User `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
}
