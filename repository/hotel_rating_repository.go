package repository

import (
	"hotel_app/entity"

	"gorm.io/gorm"
)

type IHotelRatingRepository interface {
	Create(rating *entity.HotelRating) error
	GetRatingByHotel(hotelId int) (float64, error)
}

type HotelRatingRepository struct {
	db *gorm.DB
}

func NewHotelRatingRepository(db *gorm.DB) IHotelRatingRepository {
	return &HotelRatingRepository{
		db: db,
	}
}

func (r *HotelRatingRepository) Create(rating *entity.HotelRating) error {
	result := r.db.Create(rating)
	return result.Error
}

func (r *HotelRatingRepository) GetRatingByHotel(hotelID int) (float64, error) {
	var avg float64
	err := r.db.Model(&entity.HotelRating{}).
		Select("COALESCE(AVG(rating), 0)").
		Where("hotel_id = ?", hotelID).
		Scan(&avg).Error
	if err != nil {
		return 0, err
	}
	return avg, nil
}
