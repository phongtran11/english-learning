package domain

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID          uint           `gorm:"type:bigserial;primaryKey" json:"id"`
	Email       string         `gorm:"type:varchar(255);uniqueIndex;not null" json:"email" validate:"required,email"`
	Password    string         `gorm:"type:varchar(255);not null" json:"-" validate:"required,min=8"`
	FirstName   string         `gorm:"type:varchar(100)" json:"firstName" validate:"required"`
	LastName    string         `gorm:"type:varchar(100)" json:"lastName" validate:"required"`
	PhoneNumber string         `gorm:"type:varchar(20)" json:"phoneNumber"`
	Birthdate   *time.Time     `gorm:"type:date" json:"birthdate"`
	CreatedAt   time.Time      `gorm:"type:timestamp with time zone;autoCreateTime" json:"createdAt"`
	UpdatedAt   time.Time      `gorm:"type:timestamp with time zone;autoUpdateTime" json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"type:timestamp with time zone;index" json:"-"`

	// Relations
	Sessions []Session `gorm:"foreignKey:UserID" json:"sessions,omitempty"`
}
