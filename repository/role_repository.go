package repository

import (
	"gorm.io/gorm"

	"hotel_app/entity"
)

type RoleRepository struct {
	db *gorm.DB
}

type IRoleRepository interface {
	GetByName(name string) (*entity.Role, error)
	GetById(id int) (*entity.Role, error)
}

func NewRoleRepository(db *gorm.DB) IRoleRepository {
	return &RoleRepository{
		db: db,
	}
}

func (r *RoleRepository) GetByName(name string) (*entity.Role, error) {
	var role entity.Role
	err := r.db.Where("name = ?", name).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *RoleRepository) GetById(id int) (*entity.Role, error) {
	var role entity.Role
	result := r.db.First(&role, id).Error
	if result != nil {
		return nil, result
	}
	return &role, nil
}
