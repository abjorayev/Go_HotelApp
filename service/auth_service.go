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
	LoginNotFoundError    = errors.New("login not found")
	RoleNotFoundError     = errors.New("role not found")
	UserNotFoundError     = errors.New("user not found")
)

type IAuthService interface {
	Register(dto UserCreateDTO) (string, error)
	Login(dto UserCreateDTO) (string, error)
	SetRole(userID int, roleID int) error
}

type AuthService struct {
	userRepository repository.IUserRepository
	roleRepository repository.IRoleRepository
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

func NewAuthService(userRepository repository.IUserRepository, roleRepository repository.IRoleRepository) IAuthService {
	return &AuthService{
		userRepository: userRepository,
		roleRepository: roleRepository,
	}
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
	role, err := s.roleRepository.GetByName("user")
	if err != nil {
		return "", err
	}
	user.RoleId = role.Id
	if err := s.userRepository.Register(user); err != nil {
		return "", err
	}

	return jwt.GenerateJwtToken(user.Id, "user")
}

func (s *AuthService) Login(dto UserCreateDTO) (string, error) {
	user, err := s.userRepository.GetByLogin(dto.Login)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", LoginNotFoundError
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(dto.Password)); err != nil {
		return "", ErrInvalidCredentials
	}
	if dto.Login == "admin" {
		return jwt.GenerateJwtToken(user.Id, "admin")
	}

	return jwt.GenerateJwtToken(user.Id, "user")
}

func (s *AuthService) SetRole(userID int, roleID int) error {
	role, err := s.roleRepository.GetById(roleID)
	if err != nil {
		return err
	}
	if role == nil {
		return RoleNotFoundError
	}
	user, err := s.userRepository.GetById(userID)
	if err != nil {
		return err
	}
	if user == nil {
		return UserNotFoundError
	}
	err = s.userRepository.SetRole(userID, roleID)
	if err != nil {
		return err
	}

	return nil
}
