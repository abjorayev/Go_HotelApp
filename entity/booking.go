package entity

import "time"

type Booking struct {
	ID            int `gorm:"primaryKey"`
	StartDate     time.Time
	EndDate       time.Time
	CountOfGuests int
	Price         int
	BookingStatus int `gorm:"column:booking_status;default:1"` // 1 - confirmed, 2 - booked, 3 - cancelled
	RoomID        int
	Room          Room `gorm:"foreignKey:RoomID"`
	UserId        int
	HotelID       int
}
