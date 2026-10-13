package service

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"hotel_app/entity"
	"hotel_app/pkg/jwt"
	"hotel_app/repository"
)

var (
	ErrLoginTaken         = errors.New("login already taken")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type IAuthService interface {
	Register(dto UserCreateDTO) (string, error)
	//Login(login, password string) (string, error)
}

type AuthService struct {
	userRepository repository.IUserRepository
}

type UserCreateDTO struct {
	Login    string `json:"login" binding:"required,min=3"`
	Password string `json:"password" binding:"required,min=6"`
}

func MapDTOtoEntity(dto UserCreateDTO, hashedPassword string) *entity.User {
	return &entity.User{
		Login:    dto.Login,
		Password: hashedPassword,
	}
}

func NewAuthService(userRepository repository.IUserRepository) IAuthService {
	return &AuthService{userRepository: userRepository}
}

func (s *AuthService) Register(dto UserCreateDTO) (string, error) {
	_, err := s.userRepository.GetByLogin(dto.Login)
	if err == nil {
		return "", ErrLoginTaken
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return "", err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(dto.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	user := MapDTOtoEntity(dto, string(hash))
	if err := s.userRepository.Register(user); err != nil {
		return "", err
	}

	return jwt.GenerateJwtToken(user.Id)
}
