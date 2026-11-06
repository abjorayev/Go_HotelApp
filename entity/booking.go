package entity

import "time"

type Booking struct {
	ID            int `gorm:"primaryKey"`
	StartDate     time.Time
	EndDate       time.Time
	CountOfGuests int
	IsActive      bool
	Price         int
	BookingStatus int `gorm:"default:1"` // 1 - free, 2 - booked, 3 - cancelled
	RoomID        int
	Room          Room `gorm:"foreignKey:RoomID"`
	UserId        int
	HotelID       int
}
