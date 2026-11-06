package service

import (
	"errors"
	"hotel_app/entity"
	"hotel_app/repository"
)

var (
	UserNotBookedInHotelError = errors.New("User has not booked in this hotel")
)

type IHotelRatingService interface {
	CreateRating(rating HotelRatingCreateDTO, userID int) error
	GetAverageRatingByHotel(hotelID int) (float64, error)
}

type HotelRatingCreateDTO struct {
	HotelID int `json:"hotel_id"`
	Rating  int `json:"rating"`
}

func ConvertToHotelRatingEntity(dto HotelRatingCreateDTO, userID int) *entity.HotelRating {
	return &entity.HotelRating{
		HotelID: dto.HotelID,
		UserID:  userID,
		Rating:  dto.Rating,
	}
}

type HotelRatingService struct {
	hotelRatingRepo repository.IHotelRatingRepository
	bookingService  IBookingService
	hotelRepo       repository.IHotelRepository
}

func NewHotelRatingService(hotelRatingRepo repository.IHotelRatingRepository, bookingService IBookingService, hotelRepo repository.IHotelRepository) IHotelRatingService {
	return &HotelRatingService{
		hotelRatingRepo: hotelRatingRepo,
		bookingService:  bookingService,
		hotelRepo:       hotelRepo,
	}
}

func (s *HotelRatingService) CreateRating(ratingDTO HotelRatingCreateDTO, userID int) error {
	// Check if the user has booked in this hotel
	hasBooked, err := s.bookingService.CheckIfUserBookedInThisHotel(userID, ratingDTO.HotelID)
	if err != nil {
		return err
	}
	if !hasBooked {
		return UserNotBookedInHotelError
	}

	hotel, err := s.hotelRepo.GetByID(ratingDTO.HotelID)
	if err != nil {
		return err
	}
	if hotel == nil {
		return HotelNotFoundError
	}
	err = s.hotelRatingRepo.Create(ConvertToHotelRatingEntity(ratingDTO, userID))
	if err != nil {
		return err
	}
	return nil
}

func (s *HotelRatingService) GetAverageRatingByHotel(hotelID int) (float64, error) {
	hotel, err := s.hotelRepo.GetByID(hotelID)
	if err != nil {
		return 0, err
	}
	if hotel == nil {
		return 0, HotelNotFoundError
	}
	avgRating, err := s.hotelRatingRepo.GetRatingByHotel(hotelID)
	if err != nil {
		return 0, err
	}
	return avgRating, nil
}
