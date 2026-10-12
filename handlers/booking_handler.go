package handlers

import (
	"hotel_app/service"

	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type BookingHandler struct {
	bookingService service.IBookingService
}

func NewBookingHandler(bookingService service.IBookingService) *BookingHandler {
	return &BookingHandler{bookingService: bookingService}
}

func (h *BookingHandler) CreateBooking(c *gin.Context) {
	var newBooking service.BookingCreateDTO
	if err := c.ShouldBindJSON(&newBooking); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := h.bookingService.CreateBooking(&newBooking)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Booking created successfully"})
}

func (h *BookingHandler) GetByBookingId(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	booking, err := h.bookingService.GetBookingByID(id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, booking)
}

func (h *BookingHandler) DeleteBooking(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	err2 := h.bookingService.DeleteBooking(id)
	if err2 != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
}

func (h *BookingHandler) GetByRoomid(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	bookings, err := h.bookingService.GetByRoomId(id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, bookings)

}
