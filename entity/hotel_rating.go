package entity

type HotelRating struct {
	ID      int `json:"id"`
	HotelID int `json:"hotel_id"`
	UserID  int `json:"user_id"`
	Rating  int `json:"rating"`
}
