package service

import (
	"hotel_app/entity"
	"hotel_app/repository"
)

type HotelCreateDTO struct {
	Name string `json:"name" binding:"required"`
	Address string `json:"address" binding:"required"`
	City    string `json:"city" binding:"required"`
	Stars   int    `json:"stars" binding:"required"`
	ManagerID int   `json:"manager_id" binding:"required"`
}

type HotelResponseDTO struct {
	ID      int   `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`
	City    string `json:"city"`
	Stars   int    `json:"stars"`
	ManagerID int   `json:"manager_id"`
}

type IHotelService interface {
	GetAllHotels() ([]HotelResponseDTO, error)
	GetHotelByID(id int) (*HotelResponseDTO, error)
	CreateHotel(dto *HotelCreateDTO) error
	UpdateHotel(id int, dto *HotelCreateDTO) error
	DeleteHotel(id int) error
	GetHotelsByCity(city string) ([]HotelResponseDTO, error)
	GetHotelsByManagerID(managerID int) ([]HotelResponseDTO, error)	
}


type HotelService struct {
	hotelRepository repository.IHotelRepository
}

func NewHotelService(hotelRepository repository.IHotelRepository) IHotelService {
	return &HotelService{
		hotelRepository: hotelRepository,
	}
}

func MapHotelToResponseDTO(hotel *entity.Hotel) HotelResponseDTO {
	return HotelResponseDTO{
		ID: hotel.ID,
		Name: hotel.Name,
		Address: hotel.Address,
		City: hotel.City,
		Stars: hotel.Stars,
		ManagerID: hotel.ManagerID,
	}
}

func MapCreateDTOToHotel(dto *HotelCreateDTO) entity.Hotel {
	return entity.Hotel{
	
		Name: dto.Name,
	Address: dto.Address,
	City: dto.City,
	Stars: dto.Stars,
	ManagerID: dto.ManagerID,
	}
}

func (s *HotelService) GetAllHotels() ([]HotelResponseDTO, error)  {
    hotels, err := s.hotelRepository.GetAll()
	if err != nil {
		return nil, err
	}
	var hotelDTOs []HotelResponseDTO
	for _, hotel := range hotels {
		hotelDTOs = append(hotelDTOs, MapHotelToResponseDTO(&hotel))
	}
	return hotelDTOs, nil
}

func (s *HotelService) GetHotelByID(id int) (*HotelResponseDTO, error) {
	hotel, err := s.hotelRepository.GetByID(id)
	if err != nil {
		return nil, err
	}
	dto := MapHotelToResponseDTO(hotel)
	return &dto, nil
}

func (s *HotelService) CreateHotel(dto *HotelCreateDTO) error {
	hotel := MapCreateDTOToHotel(dto)
	err := s.hotelRepository.Create(&hotel)
	if err != nil {
		return err
	}
	return nil
}

func (s *HotelService) UpdateHotel(id int, dto *HotelCreateDTO) error {
	hotel := MapCreateDTOToHotel(dto)
	hotel.ID = id
	err := s.hotelRepository.Update(&hotel)
	if err != nil {
		return err
	}
	return nil
}

func (s *HotelService) DeleteHotel(id int) error {
	err := s.hotelRepository.Delete(id)
	if err != nil {
		return err
	}
	return nil
}

func (s *HotelService) GetHotelsByCity(city string) ([]HotelResponseDTO, error) {
	hotels, err := s.hotelRepository.GetByCity(city)
	hotelDTOs := []HotelResponseDTO{}
	if err != nil {
		return nil, err
	}
	for _, hotel := range hotels {
		hotelDTOs = append(hotelDTOs, MapHotelToResponseDTO(&hotel))
	}
	return hotelDTOs, nil
}

func (s *HotelService) GetHotelsByManagerID(managerID int) ([]HotelResponseDTO, error) {
	hotels, err := s.hotelRepository.GetByManagerID(managerID)
	hotelDTOs := []HotelResponseDTO{}
	if err != nil {
		return nil, err
	}
	for _, hotel := range hotels {
		hotelDTOs = append(hotelDTOs, MapHotelToResponseDTO(&hotel))
	}
	return hotelDTOs, nil
}
