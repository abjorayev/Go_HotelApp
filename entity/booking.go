package entity

import "time"

type Booking struct {
	ID        uint   `gorm:"primaryKey"`
	StartDate time.Time
	EndDate   time.Time
	CountOfGuests int
	RoomID    uint
	Room      Room `gorm:"foreignKey:RoomID"`
}