package handlers

import (
	"hotel_app/middleware"
	"hotel_app/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type HotelRatingHandler struct {
	hotelRatingService service.IHotelRatingService
}

func NewHotelRatingHandler(service service.IHotelRatingService) *HotelRatingHandler {
	return &HotelRatingHandler{hotelRatingService: service}
}

func (h *HotelRatingHandler) CreateRating(c *gin.Context) {
	var newRating service.HotelRatingCreateDTO

	if err := c.ShouldBindJSON(&newRating); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID, ok := middleware.CurrentUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	err := h.hotelRatingService.CreateRating(newRating, userID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Rating created successfully"})
}

func (h *HotelRatingHandler) GetAverageRatingByHotel(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	avgRating, err := h.hotelRatingService.GetAverageRatingByHotel(id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"average_rating": avgRating})
}
