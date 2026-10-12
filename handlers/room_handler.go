package handlers

import (
	"hotel_app/service"

	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type RoomHandler struct {
	roomService service.IRoomService
}

func NewRoomHandler(service service.IRoomService) *RoomHandler {
	return &RoomHandler{roomService: service}
}

func (h *RoomHandler) GetByHotelId(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	rooms, err := h.roomService.GetByHotelID(id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, rooms)
}

func (h *RoomHandler) GetRoomById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	room, err := h.roomService.GetById(id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if room == nil {
		c.JSON(404, nil)
	}
	c.JSON(200, room)
}

func (h *RoomHandler) CreateRoom(c *gin.Context) {
	var newRoom service.RoomCreateDTO

	if err := c.ShouldBindJSON(&newRoom); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := h.roomService.CreateRoom(&newRoom); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, gin.H{"message": "Room created successfully"})
}

func (h *RoomHandler) UpdateRoom(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var roomUpdate service.RoomCreateDTO
	if err := c.ShouldBindJSON(&roomUpdate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.roomService.UpdateRoom(id, &roomUpdate); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Room updated successfully"})
}

func (h *RoomHandler) DeleteRoom(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.roomService.DeleteRoom(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Room deleted successfully"})
}
