package repository

import (
	"gorm.io/gorm"

	"hotel_app/entity"
)

type UserRepository struct {
	db *gorm.DB
}

type IUserRepository interface {
	Register(user *entity.User) error
	GetByLogin(login string) (*entity.User, error)
}

func NewUserRepository(db *gorm.DB) IUserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) Register(user *entity.User) error {
	result := r.db.Create(user)
	return result.Error
}

func (r *UserRepository) GetByLogin(login string) (*entity.User, error) {
	var user entity.User
	if err := r.db.Where("login = ?", login).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
