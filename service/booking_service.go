package service

import (
	"fmt"
	"hotel_app/entity"
	"hotel_app/repository"
	"time"
)

type BookingService struct {
	bookingRepo *repository.BookingRepository
	roomRepo    *repository.RoomRepository
}

func NewBookingService(bookingRepo *repository.BookingRepository, roomRepo *repository.RoomRepository) *BookingService {
	return &BookingService{
		bookingRepo: bookingRepo,
		roomRepo:    roomRepo,
	}
}

type BookingCreateDTO struct {
	RoomID    int   `json:"room_id"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	CountOfGuests int    `json:"count_of_guests"`
}

type BookingResponseDTO struct {
	ID        int   `json:"id"`
	RoomID    int   `json:"room_id"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	CountOfGuests int    `json:"count_of_guests"`
}

func ConvertToBookingResponseDTO(booking *entity.Booking) *BookingResponseDTO {
	return &BookingResponseDTO{
		ID:        booking.ID,	
	RoomID:    booking.RoomID,
		StartDate: booking.StartDate,
		EndDate:   booking.EndDate,
		CountOfGuests: booking.CountOfGuests,
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

func (s *BookingService) CreateBooking(dto *BookingCreateDTO)  error {
	// Check if the room exists
	room, err := s.roomRepo.GetByID(dto.RoomID)
	if err != nil {
		return err
	}
    if room == nil {
		return fmt.Errorf("room with ID %d not found", dto.RoomID)
	}

	if room.CountOfGuests < dto.CountOfGuests {
		return fmt.Errorf("room with ID %d cannot accommodate %d guests", dto.RoomID, dto.CountOfGuests)
	}

	aviable, err := s.bookingRepo.CheckAvailability(dto.RoomID, dto.StartDate, dto.EndDate)
	if !aviable {
		return fmt.Errorf("room with ID %d is not available from %s to %s", dto.RoomID, dto.StartDate, dto.EndDate)
	}
	if err != nil {
		return fmt.Errorf("error checking room availability: %w", err)
	}

	booking := ConvertToBookingEntity(dto)
	err = s.bookingRepo.Create(booking)
	if err != nil {
		return fmt.Errorf("error creating booking: %w", err)
	} 
   return nil
}

func (s *BookingService) GetBookingByID(id int) (*BookingResponseDTO, error) {
	booking, err := s.bookingRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	return ConvertToBookingResponseDTO(booking), nil
}

func (s *BookingService) DeleteBooking(id int) error {
	err := s.bookingRepo.Delete(id)	
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


