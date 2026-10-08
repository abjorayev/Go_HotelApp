package handlers

import (
	"hotel_app/service"

	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type HotelHandler struct {
	hotelService service.IHotelService
}

func NewHotelHandler(hotelService service.IHotelService) *HotelHandler {
	return &HotelHandler{hotelService: hotelService}
}

func (h * HotelHandler) GetHotels(c *gin.Context) {
	hotels, err := h.hotelService.GetAllHotels()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, hotels)
}

func (h *HotelHandler) GetHotelByID(c *gin.Context) {
   id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
        return
    }
	hotel, err := h.hotelService.GetHotelByID(id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, hotel)
}

func (h *HotelHandler) CreateHotel(c *gin.Context) {
	var hotelDTO service.HotelCreateDTO
	if err := c.ShouldBindJSON(&hotelDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := h.hotelService.CreateHotel(&hotelDTO)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(201, gin.H{"message": "Hotel created successfully"})
}

func (h *HotelHandler) UpdateHotel(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var hotelDTO service.HotelCreateDTO
	if err := c.ShouldBindJSON(&hotelDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err2 := h.hotelService.UpdateHotel(id, &hotelDTO)
	if err2 != nil {
		c.JSON(500, gin.H{"error": err2.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Hotel updated successfully"})
}

func (h *HotelHandler) DeleteHotel(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	err2 := h.hotelService.DeleteHotel(id)
	if err2 != nil {
		c.JSON(500, gin.H{"error": err2.Error()})
		return
	}	
}