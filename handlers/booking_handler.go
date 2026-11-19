package handlers

import (
	"hotel_app/service"

	"hotel_app/middleware"
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
	userID, ok := middleware.CurrentUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	err := h.bookingService.CreateBooking(&newBooking, userID)
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
	userId, ok := middleware.CurrentUserID(c)
	role, ok := middleware.CurrentUserRole(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if err := h.bookingService.DeleteBooking(id, userId, role); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "booking deleted successfully"})
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

func (h *BookingHandler) GetByUserId(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	pageStr := c.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid page number"})
		return
	}

	rowsStr := c.DefaultQuery("rows", "10")
	rows, err := strconv.Atoi(rowsStr)
	if err != nil || rows <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid rows count"})
		return
	}
	bookings, err := h.bookingService.GetByUserId(id, page, rows)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, bookings)
}
