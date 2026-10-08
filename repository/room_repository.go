package repository

import (
	"gorm.io/gorm"

"hotel_app/entity")

type RoomRepository struct {
	db *gorm.DB
}

func NewRoomRepository(db *gorm.DB) *RoomRepository {
	return &RoomRepository{
		db: db,
	}
}

func (r *RoomRepository) GetAll() ([]entity.Room, error) {
	var rooms []entity.Room
	result := r.db.Find(&rooms)
	return rooms, result.Error
}

func (r *RoomRepository) GetByID(id int) (*entity.Room, error) {
	var room entity.Room
	result := r.db.First(&room, id)
	return &room, result.Error
}

func (r *RoomRepository) Create(room *entity.Room) error {
	result := r.db.Create(room)
	return result.Error
}

func (r *RoomRepository) Update(room *entity.Room) error {
	result := r.db.Save(room)
	return result.Error
}

func (r *RoomRepository) Delete(id int) error {
	result := r.db.Delete(&entity.Room{}, id)
	return result.Error
}

func (r *RoomRepository) GetByHotelID(hotelID int) ([]entity.Room, error) {
	var rooms []entity.Room
	result := r.db.Where("hotel_id = ?", hotelID).Find(&rooms)
	return rooms, result.Error
}