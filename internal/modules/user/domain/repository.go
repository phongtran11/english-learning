package domain

import "errors"

var ErrUserNotFound = errors.New("user not found")

type UserRepository interface {
	Create(user *User) error
	FindByEmail(email string) (*User, error)
	FindByID(id uint) (*User, error)
	Update(user *User) error
	Delete(id uint) error
	List(offset, limit int) ([]User, int64, error)
}
