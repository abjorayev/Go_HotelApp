package repository

import (
	"hotel_app/entity"
	"time"

	"gorm.io/gorm"
)

type BookingRepository struct {
	db *gorm.DB
}

type IBookingRepository interface {
	Create(booking *entity.Booking) error
	GetByID(id int) (*entity.Booking, error)
	GetByRoomID(roomID int) ([]entity.Booking, error)
	Update(booking *entity.Booking) error
	Delete(id int) error
	CheckAvailability(roomID int, startDate, endDate time.Time) (bool, error)
	GetByUserId(userId int) ([]entity.Booking, error)
	CheckIfUserBookedInThisHotel(userId int, hotelId int) (bool, error)
}

func NewBookingRepository(db *gorm.DB) IBookingRepository {
	return &BookingRepository{
		db: db,
	}
}

func (r *BookingRepository) Create(booking *entity.Booking) error {
	result := r.db.Create(booking)
	return result.Error
}

func (r *BookingRepository) GetByID(id int) (*entity.Booking, error) {
	var booking entity.Booking
	result := r.db.First(&booking, id)
	return &booking, result.Error
}

func (r *BookingRepository) GetByRoomID(roomID int) ([]entity.Booking, error) {
	var bookings []entity.Booking
	result := r.db.Where("room_id = ?", roomID).Find(&bookings)
	return bookings, result.Error
}

func (r *BookingRepository) GetByUserId(userId int) ([]entity.Booking, error) {
	var bookings []entity.Booking
	result := r.db.Where("user_id = ?", userId).Find(&bookings)
	return bookings, result.Error
}

func (r *BookingRepository) Update(booking *entity.Booking) error {
	result := r.db.Save(booking)
	return result.Error
}

func (r *BookingRepository) Delete(id int) error {
	result := r.db.Delete(&entity.Booking{}, id)
	return result.Error
}

func (r *BookingRepository) CheckAvailability(
	roomID int,
	startDate, endDate time.Time,
) (bool, error) {

	var count int64

	err := r.db.Raw(`
		SELECT COUNT(*)
		FROM bookings
		WHERE room_id = ?
		  AND start_date < ?
		  AND end_date > ?
		  AND status = 2
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

func (r *BookingRepository) CheckIfUserBookedInThisHotel(userId int, hotelId int) (bool, error) {
	var count int64
	err := r.db.Raw(`
		SELECT COUNT(*)
		FROM bookings
		WHERE user_id = ?
		  AND hotel_id = ?
		  AND status = 2
	`, userId, hotelId).Scan(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}
