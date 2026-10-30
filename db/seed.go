package db

import (
	"errors"
	"hotel_app/entity"
	"os"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Seed(db *gorm.DB) error {
	roles := []entity.Role{
		{Id: 1, Name: "user"},
		{Id: 2, Name: "admin"},
	}

	var adminRole entity.Role
	for _, r := range roles {
		var role entity.Role
		err := db.Where(entity.Role{Name: r.Name}).
			Attrs(entity.Role{Id: r.Id}).
			FirstOrCreate(&role).Error
		if err != nil {
			return err
		}
		if r.Name == "admin" {
			adminRole = role
		}
	}

	login := os.Getenv("ADMIN_LOGIN")
	password := os.Getenv("ADMIN_PASSWORD")
	if login == "" && password == "" {
		return errors.New("Login and password are null")
	}
	var count int64
	if err := db.Model(&entity.User{}).Where("login = ?", login).Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return nil
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	admin := entity.User{
		Login:    login,
		Password: string(hash),
		RoleId:   adminRole.Id,
	}
	return db.Create(&admin).Error
}
