package handlers

import (
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
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}
