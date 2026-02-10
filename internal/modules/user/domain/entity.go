package domain

import (
	"time"
)

type User struct {
	ID          uint
	Email       string
	Password    string
	FirstName   string
	LastName    string
	PhoneNumber string
	Birthdate   *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
	// DeletedAt removed as it's persistence concern, or Changed to *time.Time if logical delete is domain concept.
	// For now, I will remove it to be strictly pure as requested, assuming logical delete is an implementation detail of persistence.
	// If domain logic requires knowing if a user is deleted, I would add `IsDeleted bool` or `DeletedAt *time.Time`.
	// Given the previous code used gorm.DeletedAt, I'll omit it for now to follow "pure" instructions strictly.

}
