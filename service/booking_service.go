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
	hotelRepo   repository.IHotelRepository
}

type IBookingService interface {
	CreateBooking(dto *BookingCreateDTO, userId int) error
	GetBookingByID(id int, userId int) (*BookingResponseDTO, error)
	DeleteBooking(id int, userId int) error
	GetByRoomId(roomID int, userId int, roleName string) ([]*BookingResponseDTO, error)
	GetByUserId(userId int, page int, rows int) ([]*BookingResponseDTO, error)
	CheckIfUserBookedInThisHotel(userId int, hotelId int) (bool, error)
}

var (
	BookingNotFoundError   = fmt.Errorf("booking not found")
	UserNotAuthorizedError = fmt.Errorf("user not authorized to cancel this booking")
	BookingIsNotYoursError = fmt.Errorf("You can see other people's booking")
	HotelNotAccessError    = fmt.Errorf("You don't have access to this hotel")
)

func NewBookingService(bookingRepo repository.IBookingRepository, roomRepo repository.IRoomRepository, hotelRepo repository.IHotelRepository) IBookingService {
	return &BookingService{
		bookingRepo: bookingRepo,
		roomRepo:    roomRepo,
		hotelRepo:   hotelRepo,
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
	hotel, err := s.hotelRepo.GetByID(room.HotelID)
	if err != nil {
		return err
	}
	if hotel.IsActive == false {
		return fmt.Errorf("Hotel is not active")
	}
	if nightsBetween(dto.StartDate, dto.EndDate) < 1 {
		return fmt.Errorf("start date cannot be after end date")
	}
	if dto.EndDate.Before(time.Now()) {
		return fmt.Errorf("time of the booking cannot be before real time")
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
	//booking.IsActive = true
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

func (s *BookingService) GetBookingByID(id int, userId int) (*BookingResponseDTO, error) {
	booking, err := s.bookingRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if booking.UserId != userId {
		return nil, BookingIsNotYoursError
	}
	return ConvertToBookingResponseDTO(booking), nil
}

func (s *BookingService) DeleteBooking(id int, userId int) error {
	booking, err := s.bookingRepo.GetByID(id)
	if err != nil {
		return err
	}
	if booking == nil {
		return BookingNotFoundError
	}
	if booking.BookingStatus == entity.BookingStatusCancelled {
		return fmt.Errorf("Booking already cancelled")
	}
	if booking.UserId != userId {
		return UserNotAuthorizedError
	}
	booking.BookingStatus = entity.BookingStatusCancelled
	//booking.IsActive = false
	err = s.bookingRepo.Update(booking)
	if err != nil {
		return fmt.Errorf("error cancelling booking: %w", err)
	}
	return nil
}

func (s *BookingService) GetByRoomId(roomID int, userId int, roleName string) ([]*BookingResponseDTO, error) {
	bookings, err := s.bookingRepo.GetByRoomID(roomID)
	if err != nil {
		return nil, err
	}
	managerId, err := s.roomRepo.GetManagerByRoomId(roomID)
	if err != nil {
		return nil, err
	}
	if managerId != userId && roleName != "admin" {
		return nil, HotelNotAccessError
	}
	var bookingDTOs []*BookingResponseDTO
	for _, booking := range bookings {
		bookingDTOs = append(bookingDTOs, ConvertToBookingResponseDTO(&booking))
	}
	return bookingDTOs, nil
}

func (s *BookingService) GetByUserId(userID int, page int, rows int) ([]*BookingResponseDTO, error) {
	bookings, err := s.bookingRepo.GetByUserId(userID, page, rows)
	if err != nil {
		return nil, err
	}
	var result []*BookingResponseDTO
	for _, booking := range bookings {
		result = append(result, ConvertToBookingResponseDTO(&booking))
	}
	return result, nil
}
