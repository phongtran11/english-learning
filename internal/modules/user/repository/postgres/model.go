package postgres

import (
	"english-learning/internal/modules/user/domain"
	"time"

	"gorm.io/gorm"
)

type UserGorm struct {
	ID          uint           `gorm:"type:bigserial;primaryKey"`
	Email       string         `gorm:"type:varchar(255);uniqueIndex;not null"`
	Password    string         `gorm:"type:varchar(255);not null"`
	FirstName   string         `gorm:"type:varchar(100)"`
	LastName    string         `gorm:"type:varchar(100)"`
	PhoneNumber string         `gorm:"type:varchar(20)"`
	Birthdate   *time.Time     `gorm:"type:date"`
	CreatedAt   time.Time      `gorm:"type:timestamp with time zone;autoCreateTime"`
	UpdatedAt   time.Time      `gorm:"type:timestamp with time zone;autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `gorm:"type:timestamp with time zone;index"`


}

func (m *UserGorm) ToDomain() *domain.User {
    if m == nil {
        return nil
    }
	return &domain.User{
		ID:          m.ID,
		Email:       m.Email,
		Password:    m.Password,
		FirstName:   m.FirstName,
		LastName:    m.LastName,
		PhoneNumber: m.PhoneNumber,
		Birthdate:   m.Birthdate,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
		// DeletedAt is not part of pure domain usually, or we can add it if needed. 
        // Plan says remove GORM tags. If domain has DeletedAt as time.Time or custom struct, we map it. 
        // Checking domain/user.go again, it has gorm.DeletedAt. I should change that to time.Time or remove it.
        // For now, I will assume domain will not have gorm.DeletedAt. 
	}
}

func FromDomainUser(u *domain.User) *UserGorm {
    if u == nil {
        return nil
    }
	return &UserGorm{
		ID:          u.ID,
		Email:       u.Email,
		Password:    u.Password,
		FirstName:   u.FirstName,
		LastName:    u.LastName,
		PhoneNumber: u.PhoneNumber,
		Birthdate:   u.Birthdate,
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
	}
}
