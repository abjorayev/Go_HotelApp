package entity

import "time"

type Room struct {
	ID        int      `json:"id"`
	HotelID   int      `json:"hotel_id"`
	RoomNumber string    `json:"room_number"`
	RoomType  string    `json:"room_type"`
	Price     float64   `json:"price"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	CountOfGuests int
	Bookings  []Booking `gorm:"foreignKey:RoomID"`
}
