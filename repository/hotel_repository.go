package repository

import ("gorm.io/gorm" 

"hotel_app/entity")

type IHotelRepository interface {
    GetAll() ([]entity.Hotel, error)
    GetByID(id int) (*entity.Hotel, error)
    Create(hotel *entity.Hotel) error
    Update(hotel *entity.Hotel) error
    Delete(id int) error
    GetByCity(city string) ([]entity.Hotel, error)
    GetByManagerID(managerID int) ([]entity.Hotel, error)
}

type HotelRepository struct {
	db *gorm.DB
}

func NewHotelRepository(db *gorm.DB) IHotelRepository {
	return &HotelRepository{
		db: db,
	}
}

func (r *HotelRepository) GetAll() ([]entity.Hotel, error) {
	var hotels []entity.	Hotel
	result := r.db.Find(&hotels)
	return hotels, result.Error
}

func (r *HotelRepository) GetByID(id int) (*entity.Hotel, error) {
	var hotel entity.Hotel
	result := r.db.First(&hotel, id)
	return &hotel, result.Error
}

func (r *HotelRepository) Create(hotel *entity.Hotel) error {
	result := r.db.Create(hotel)
	return result.Error
}

func (r *HotelRepository) Update(hotel *entity.Hotel) error {
	result := r.db.Save(hotel)
	return result.Error
}

func (r *HotelRepository) Delete(id int) error {
	result := r.db.Delete(&entity.Hotel{}, id)
	return result.Error
}

func (r* HotelRepository) GetByCity(city string) ([]entity.Hotel, error) {
	var hotels []entity.Hotel
	result := r.db.Where("city = ?", city).Find(&hotels)
	return hotels, result.Error
}

func (r *HotelRepository) GetByManagerID(managerID int) ([]entity.Hotel, error) {
	var hotels []entity.Hotel
	result := r.db.Where("manager_id = ?", managerID).Find(&hotels)
	return hotels, result.Error
}