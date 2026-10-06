package service

import (
	"hotel_app/repository"

	"hotel_app/entity"
)


type RoomService struct {
	roomRepository repository.RoomRepository
}

func NewRoomService(roomRepository repository.RoomRepository) *RoomService {
	return &RoomService{
		roomRepository: roomRepository,
	}
}

type RoomCreateDTO struct {
	HotelID    uint    `json:"hotel_id" binding:"required"`
	RoomNumber string  `json:"room_number" binding:"required"`
	RoomType   string  `json:"room_type" binding:"required"`
	Price      float64 `json:"price" binding:"required"`
	IsActive   bool    `json:"is_active" binding:"required"`
}

type RoomResponseDTO struct {
	ID        uint    `json:"id"`
	HotelID   uint    `json:"hotel_id"`
	RoomNumber string  `json:"room_number"`
	RoomType  string  `json:"room_type"`
	Price     float64 `json:"price"`
	IsActive  bool    `json:"is_active"`
}

func MapCreateDTOToRoom(dto *RoomCreateDTO) entity.Room {
	return entity.Room{
		HotelID:    dto.HotelID,
		RoomNumber: dto.RoomNumber,
		RoomType:   dto.RoomType,
		Price:      dto.Price,
		IsActive:   dto.IsActive,
	}
}

func MapRoomToResponseDTO(room *entity.Room) RoomResponseDTO {
	return RoomResponseDTO{
		ID:        room.ID,	
	HotelID:   room.HotelID,
	RoomNumber: room.RoomNumber,
	RoomType:  room.RoomType,
	Price:     room.Price,
	IsActive:  room.IsActive,
	}
}

func (s *RoomService) CreateRoom(dto *RoomCreateDTO) error {
	room := MapCreateDTOToRoom(dto)
	return s.roomRepository.Create(&room)
}

func (s *RoomService) UpdateRoom(id uint, dto *RoomCreateDTO) error {
	room, err := s.roomRepository.GetByID(id)
	if err != nil {
		return err
	}
    roomEntity := MapCreateDTOToRoom(dto)
	roomEntity.ID = room.ID
	return s.roomRepository.Update(&roomEntity)
}

func (s *RoomService) DeleteRoom(id uint) error {
	return s.roomRepository.Delete(id)
}

func (s *RoomService) GetById(id uint) (*RoomResponseDTO, error) {
	room, err := s.roomRepository.GetByID(id)
	if err != nil {
		return nil, err
	}
	dto := MapRoomToResponseDTO(room)
	return &dto, nil
}

func (s *RoomService) GetByHotelID(hotelID uint) ([]RoomResponseDTO, error) {
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