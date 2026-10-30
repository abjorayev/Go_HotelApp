package entity

import "time"

type Booking struct {
	ID            int `gorm:"primaryKey"`
	StartDate     time.Time
	EndDate       time.Time
	CountOfGuests int
	IsActive      bool
	RoomID        int
	Room          Room `gorm:"foreignKey:RoomID"`
	UserId        int
}
