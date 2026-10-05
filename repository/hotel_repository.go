package repository

import ("gorm.io/gorm" 

"hotel_app/entity")

type HotelRepository struct {
	db *gorm.DB
}

func NewHotelRepository(db *gorm.DB) *HotelRepository {
	return &HotelRepository{
		db: db,
	}
}

func (r *HotelRepository) GetAll() ([]entity.Hotel, error) {
	var hotels []entity.	Hotel
	result := r.db.Find(&hotels)
	return hotels, result.Error
}

func (r *HotelRepository) GetByID(id uint) (*entity.Hotel, error) {
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

func (r *HotelRepository) Delete(id uint) error {
	result := r.db.Delete(&entity.Hotel{}, id)
	return result.Error
}