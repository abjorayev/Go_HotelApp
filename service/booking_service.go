package service

import (
	"fmt"
	"hotel_app/entity"
	"hotel_app/repository"
	"time"
)

type BookingService struct {
	bookingRepo repository.IBookingRepository
	roomRepo    repository.IRoomRepository
}

type IBookingService interface {
	CreateBooking(dto *BookingCreateDTO, userId int) error
	GetBookingByID(id int) (*BookingResponseDTO, error)
	DeleteBooking(id int) error
	GetByRoomId(roomID int) ([]*BookingResponseDTO, error)
	GetByUserId(userId int) ([]*BookingResponseDTO, error)
	CancelBooking(id int, userId int) error
	CheckIfUserBookedInThisHotel(userId int, hotelId int) (bool, error)
}

var (
	BookingNotFoundError   = fmt.Errorf("booking not found")
	UserNotAuthorizedError = fmt.Errorf("user not authorized to cancel this booking")
)

func NewBookingService(bookingRepo repository.IBookingRepository, roomRepo repository.IRoomRepository) IBookingService {
	return &BookingService{
		bookingRepo: bookingRepo,
		roomRepo:    roomRepo,
	}
}

type BookingCreateDTO struct {
	RoomID        int       `json:"room_id"`
	StartDate     time.Time `json:"start_date"`
	EndDate       time.Time `json:"end_date"`
	CountOfGuests int       `json:"count_of_guests"`
}

type BookingResponseDTO struct {
	ID            int       `json:"id"`
	RoomID        int       `json:"room_id"`
	StartDate     time.Time `json:"start_date"`
	EndDate       time.Time `json:"end_date"`
	CountOfGuests int       `json:"count_of_guests"`
	Price         int       `json:"price"`
}

func ConvertToBookingResponseDTO(booking *entity.Booking) *BookingResponseDTO {
	return &BookingResponseDTO{
		ID:            booking.ID,
		RoomID:        booking.RoomID,
		StartDate:     booking.StartDate,
		EndDate:       booking.EndDate,
		CountOfGuests: booking.CountOfGuests,
		Price:         booking.Price,
	}
}

func ConvertToBookingEntity(dto *BookingCreateDTO) *entity.Booking {
	return &entity.Booking{
		RoomID:        dto.RoomID,
		StartDate:     dto.StartDate,
		EndDate:       dto.EndDate,
		CountOfGuests: dto.CountOfGuests,
	}
}

func (s *BookingService) CheckIfUserBookedInThisHotel(userId int, hotelId int) (bool, error) {
	return s.bookingRepo.CheckIfUserBookedInThisHotel(userId, hotelId)
}

func (s *BookingService) CreateBooking(dto *BookingCreateDTO, userId int) error {
	// Check if the room exists
	room, err := s.roomRepo.GetByID(dto.RoomID)
	if err != nil {
		return err
	}
	if room == nil {
		return fmt.Errorf("room with ID %d not found", dto.RoomID)
	}

	if dto.StartDate.After(dto.EndDate) {
		return fmt.Errorf("start date cannot be after end date")
	}
	if dto.CountOfGuests <= 0 {
		return fmt.Errorf("count of guests must be greater than zero")
	}

	if room.IsActive == false {
		return fmt.Errorf("room with ID %d is not active", dto.RoomID)
	}

	if room.CountOfGuests < dto.CountOfGuests {
		return fmt.Errorf("room with ID %d cannot accommodate %d guests", dto.RoomID, dto.CountOfGuests)
	}

	aviable, err := s.bookingRepo.CheckAvailability(dto.RoomID, dto.StartDate, dto.EndDate)
	if err != nil {
		return fmt.Errorf("error checking room availability: %w", err)
	}

	if !aviable {
		return fmt.Errorf("room with ID %d is not available from %s to %s", dto.RoomID, dto.StartDate, dto.EndDate)
	}
	days := nightsBetween(dto.StartDate, dto.EndDate)
	price := int(float64(days) * room.Price)
	booking := ConvertToBookingEntity(dto)
	booking.IsActive = true
	booking.UserId = userId
	booking.Price = price
	booking.BookingStatus = entity.BookingStatusBooked
	booking.HotelID = room.HotelID
	err = s.bookingRepo.Create(booking)
	if err != nil {
		return fmt.Errorf("error creating booking: %w", err)
	}
	return nil
}

func nightsBetween(start, end time.Time) int {
	s := time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, time.UTC)
	e := time.Date(end.Year(), end.Month(), end.Day(), 0, 0, 0, 0, time.UTC)
	return int(e.Sub(s).Hours() / 24)
}

func (s *BookingService) GetBookingByID(id int) (*BookingResponseDTO, error) {
	booking, err := s.bookingRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	return ConvertToBookingResponseDTO(booking), nil
}

func (s *BookingService) DeleteBooking(id int) error {
	booking, err := s.bookingRepo.GetByID(id)
	booking.IsActive = false
	err = s.bookingRepo.Update(booking)
	//err := s.bookingRepo.Delete(id)
	if err != nil {
		return fmt.Errorf("error deleting booking: %w", err)
	}
	return nil
}

func (s *BookingService) GetByRoomId(roomID int) ([]*BookingResponseDTO, error) {
	bookings, err := s.bookingRepo.GetByRoomID(roomID)
	if err != nil {
		return nil, err
	}
	var bookingDTOs []*BookingResponseDTO
	for _, booking := range bookings {
		bookingDTOs = append(bookingDTOs, ConvertToBookingResponseDTO(&booking))
	}
	return bookingDTOs, nil
}

func (s *BookingService) GetByUserId(userID int) ([]*BookingResponseDTO, error) {
	bookings, err := s.bookingRepo.GetByUserId(userID)
	if err != nil {
		return nil, err
	}
	var result []*BookingResponseDTO
	for _, booking := range bookings {
		result = append(result, ConvertToBookingResponseDTO(&booking))
	}
	return result, nil
}

func (s *BookingService) CancelBooking(id int, userId int) error {
	booking, err := s.bookingRepo.GetByID(id)
	if err != nil {
		return err
	}
	if booking == nil {
		return BookingNotFoundError
	}
	if booking.UserId != userId {
		return UserNotAuthorizedError
	}
	booking.BookingStatus = entity.BookingStatusCancelled
	err = s.bookingRepo.Update(booking)
	if err != nil {
		return fmt.Errorf("error cancelling booking: %w", err)
	}
	return nil
}
