package http

import (
	"english-learning/internal/modules/user/domain"
	"time"
)

type RegisterRequestDTO struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type UserResponseDTO struct {
	ID          uint       `json:"id"`
	Email       string     `json:"email"`
	FirstName   string     `json:"firstName"`
	LastName    string     `json:"lastName"`
	PhoneNumber string     `json:"phoneNumber"`
	Birthdate   *time.Time `json:"birthdate"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
}

type UpdateUserRequestDTO struct {
	FirstName   string     `json:"firstName" binding:"required"`
	LastName    string     `json:"lastName" binding:"required"`
	PhoneNumber string     `json:"phoneNumber"`
	Birthdate   *time.Time `json:"birthdate"`
}

func ToUserResponse(user *domain.User) UserResponseDTO {
	if user == nil {
		return UserResponseDTO{}
	}
	return UserResponseDTO{
		ID:          user.ID,
		Email:       user.Email,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		PhoneNumber: user.PhoneNumber,
		Birthdate:   user.Birthdate,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}
}

func ToUserListResponse(users []domain.User) []UserResponseDTO {
	dtos := make([]UserResponseDTO, len(users))
	for i, user := range users {
		dtos[i] = ToUserResponse(&user)
	}
	return dtos
}
