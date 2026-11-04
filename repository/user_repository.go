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
	SetRole(userID int, roleID int) error
	GetById(id int) (*entity.User, error)
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

func (r *UserRepository) GetById(id int) (*entity.User, error) {
	var user entity.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) SetRole(userID int, roleID int) error {
	var user entity.User
	result := r.db.First(&user, userID)
	if result.Error != nil {
		return result.Error
	}
	user.RoleId = roleID
	return r.db.Save(&user).Error
}
