package postgres

import (
	"english-learning/internal/modules/user/domain"
	"errors"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) domain.UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *domain.User) error {
	userModel := FromDomainUser(user)
	if err := r.db.Create(userModel).Error; err != nil {
		return err
	}
	// Update ID back to domain
	user.ID = userModel.ID
    user.CreatedAt = userModel.CreatedAt
    user.UpdatedAt = userModel.UpdatedAt
	return nil
}

func (r *UserRepository) FindByEmail(email string) (*domain.User, error) {
	var userModel User
	err := r.db.Where("email = ?", email).First(&userModel).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}
	return userModel.ToDomain(), nil
}

func (r *UserRepository) FindByID(id uint) (*domain.User, error) {
	var userModel User
	err := r.db.First(&userModel, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}
	return userModel.ToDomain(), nil
}

func (r *UserRepository) Update(user *domain.User) error {
	userModel := FromDomainUser(user)
	return r.db.Save(userModel).Error
}

func (r *UserRepository) Delete(id uint) error {
	return r.db.Delete(&User{}, id).Error
}

func (r *UserRepository) List(offset, limit int) ([]domain.User, int64, error) {
	var userModels []User
	var count int64

	if err := r.db.Model(&User{}).Count(&count).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.Offset(offset).Limit(limit).Find(&userModels).Error
	if err != nil {
		return nil, 0, err
	}

	users := make([]domain.User, len(userModels))
	for i, model := range userModels {
		users[i] = *model.ToDomain()
	}

	return users, count, nil
}
