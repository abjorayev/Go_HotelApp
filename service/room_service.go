package service

import (
	"fmt"
	"hotel_app/repository"

	"hotel_app/entity"
)

type RoomService struct {
	roomRepository  repository.IRoomRepository
	hotelRepository repository.IHotelRepository
}

func NewRoomService(roomRepository repository.IRoomRepository, hotelRepository repository.IHotelRepository) IRoomService {
	return &RoomService{
		roomRepository:  roomRepository,
		hotelRepository: hotelRepository,
	}
}

var (
	HotelNotFoundError = fmt.Errorf("hotel not found")
)

type RoomCreateDTO struct {
	HotelID       int     `json:"hotel_id" binding:"required"`
	RoomNumber    string  `json:"room_number" binding:"required"`
	RoomType      string  `json:"room_type" binding:"required"`
	Price         float64 `json:"price" binding:"required"`
	IsActive      bool    `json:"is_active" binding:"required"`
	CountOfGuests int     `json:"count_of_guests" binding:"required"`
}

type RoomResponseDTO struct {
	ID            int     `json:"id"`
	HotelID       int     `json:"hotel_id"`
	RoomNumber    string  `json:"room_number"`
	RoomType      string  `json:"room_type"`
	Price         float64 `json:"price"`
	IsActive      bool    `json:"is_active"`
	CountOfGuests int     `json:"count_of_guests" binding:"required"`
}

type IRoomService interface {
	CreateRoom(dto *RoomCreateDTO) error
	UpdateRoom(id int, dto *RoomCreateDTO, managerId int, role string) error
	DeleteRoom(id int, managerId int, role string) error
	GetById(id int) (*RoomResponseDTO, error)
	GetByHotelID(hotelID int) ([]RoomResponseDTO, error)
}

func MapCreateDTOToRoom(dto *RoomCreateDTO) entity.Room {
	return entity.Room{
		HotelID:       dto.HotelID,
		RoomNumber:    dto.RoomNumber,
		RoomType:      dto.RoomType,
		Price:         dto.Price,
		CountOfGuests: dto.CountOfGuests,
		IsActive:      dto.IsActive,
	}
}

func MapRoomToResponseDTO(room *entity.Room) RoomResponseDTO {
	return RoomResponseDTO{
		ID:            room.ID,
		HotelID:       room.HotelID,
		RoomNumber:    room.RoomNumber,
		RoomType:      room.RoomType,
		Price:         room.Price,
		IsActive:      room.IsActive,
		CountOfGuests: room.CountOfGuests,
	}
}

func (s *RoomService) IsCurrentHotelManager(userId int, hotelId int, role string) (bool, error) {
	hotel, err := s.hotelRepository.GetByID(hotelId)
	if err != nil {
		return false, err
	}
	if hotel == nil {
		return false, HotelNotFoundError
	}
	if role == "admin" {
		return true, nil
	}
	return hotel.ManagerID == userId, nil
}

func (s *RoomService) CreateRoom(dto *RoomCreateDTO) error {
	hotel, err := s.hotelRepository.GetByID(dto.HotelID)
	if err != nil {
		return err
	}
	if hotel == nil {
		return HotelNotFoundError
	}

	room := MapCreateDTOToRoom(dto)
	return s.roomRepository.Create(&room)
}

func (s *RoomService) UpdateRoom(id int, dto *RoomCreateDTO, managerId int, role string) error {
	room, err := s.roomRepository.GetByID(id)
	if err != nil {
		return err
	}
	isManager, err := s.IsCurrentHotelManager(managerId, room.HotelID, role)
	if err != nil {
		return err
	}
	if !isManager {
		return fmt.Errorf("user is not the manager of the hotel")
	}
	//roomEntity := MapCreateDTOToRoom(dto)
	//roomEntity.ID = room.ID
	room.HotelID = dto.HotelID
	room.RoomNumber = dto.RoomNumber
	room.RoomType = dto.RoomType
	room.Price = dto.Price
	room.IsActive = dto.IsActive
	room.CountOfGuests = dto.CountOfGuests
	return s.roomRepository.Update(room)
}

func (s *RoomService) DeleteRoom(id int, managerId int, role string) error {
	room, err := s.roomRepository.GetByID(id)
	if err != nil {
		return err
	}
	isManager, err := s.IsCurrentHotelManager(managerId, room.HotelID, role)
	if err != nil {
		return err
	}
	if !isManager {
		return fmt.Errorf("user is not the manager of the hotel")
	}
	room.IsActive = false
	return s.roomRepository.Update(room)
}

func (s *RoomService) GetById(id int) (*RoomResponseDTO, error) {
	room, err := s.roomRepository.GetByID(id)
	if err != nil {
		return nil, err
	}
	dto := MapRoomToResponseDTO(room)
	return &dto, nil
}

func (s *RoomService) GetByHotelID(hotelID int) ([]RoomResponseDTO, error) {
	rooms, err := s.roomRepository.GetByHotelID(hotelID)
	if err != nil {
		return nil, err
	}
	var dtos []RoomResponseDTO
	for _, room := range rooms {
		dto := MapRoomToResponseDTO(&room)
		dtos = append(dtos, dto)
	}
	return dtos, nil
}
