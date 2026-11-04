package handlers

import (
	"errors"
	"hotel_app/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService service.IAuthService
}

func NewAuthHandler(authService service.IAuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (s *AuthHandler) Register(c *gin.Context) {
	var user service.UserCreateDTO
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	token, err := s.authService.Register(user)

	if err != nil {
		if errors.Is(err, service.ErrLoginTaken) {
			c.JSON(409, gin.H{"error": "login already taken"})
			return
		}
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (s *AuthHandler) Login(c *gin.Context) {
	var user service.UserCreateDTO
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	token, err := s.authService.Login(user)
	if err != nil {
		if errors.Is(err, service.LoginNotFoundError) {
			c.JSON(404, gin.H{"error": err.Error()})
			return
		}
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (s *AuthHandler) SetRole(c *gin.Context) {
	var request struct {
		UserID int `json:"user_id" binding:"required"`
		RoleID int `json:"role_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := s.authService.SetRole(request.UserID, request.RoleID)
	if err != nil {
		if errors.Is(err, service.UserNotFoundError) {
			c.JSON(404, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, service.RoleNotFoundError) {
			c.JSON(404, gin.H{"error": err.Error()})
			return
		}
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "role set successfully"})
}
