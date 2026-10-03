package db

import (
	"hotel_app/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(config *config.Config) (*gorm.DB, error) {
  dsn := "host=" + config.DBHost + " user=" + config.DBUser + " dbname=" + config.DBName + " password=" + config.DBPassword + " sslmode=disable"
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}