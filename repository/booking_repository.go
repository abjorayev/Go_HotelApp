package repository

import (
	"gorm.io/gorm"
	"hotel_app/entity"
    "time"
)

type BookingRepository struct {
	db *gorm.DB
}

func NewBookingRepository(db *gorm.DB) *BookingRepository {
	return &BookingRepository{
		db: db,
	}
}

func (r *BookingRepository) Create(booking *entity.Booking) error {
	result := r.db.Create(booking)
	return result.Error
}

func (r *BookingRepository) GetByID(id uint) (*entity.Booking, error) {
	var booking entity.Booking
	result := r.db.First(&booking, id)
	return &booking, result.Error
}

func (r *BookingRepository) GetByRoomID(roomID uint) ([]entity.Booking, error) {
	var bookings []entity.Booking
	result := r.db.Where("room_id = ?", roomID).Find(&bookings)
	return bookings, result.Error
}

func (r *BookingRepository) Update(booking *entity.Booking) error {
	result := r.db.Save(booking)
	return result.Error
}

func (r *BookingRepository) Delete(id uint) error {
	result := r.db.Delete(&entity.Booking{}, id)
	return result.Error
}

func (r *BookingRepository) CheckAvailability(
	roomID uint,
	startDate, endDate time.Time,
) (bool, error) {

	var count int64

	err := r.db.Raw(`
		SELECT COUNT(*)
		FROM bookings
		WHERE room_id = ?
		  AND start_date < ?
		  AND end_date > ?
	`,
		roomID,
		endDate,
		startDate,
	).Scan(&count).Error

	if err != nil {
		return false, err
	}

	return count == 0, nil
}

