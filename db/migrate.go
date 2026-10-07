package db

import (
	"gorm.io/gorm"

	"hotel_app/entity"
)

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&entity.Hotel{},
		&entity.Room{},
		&entity.Booking{},
	)
}