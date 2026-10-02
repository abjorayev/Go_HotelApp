package db


import ( "gorm.io/gorm"
"gorm.io/driver/postgres"
"hotel_app/config"
)

func main(config *config.Config) (*gorm.DB, error) {
  dsn := "host=" + config.DBHost + " user=" + config.DBUser + " dbname=" + config.DBName + " password=" + config.DBPassword + " sslmode=disable"
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}