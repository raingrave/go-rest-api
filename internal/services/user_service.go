package services

import (
	"github.com/google/uuid"
	"github.com/raingrave/apirest/internal/models"
	"github.com/raingrave/apirest/internal/repositories"
)

type UserService interface {
	CreateUser(user models.User) (models.User, error)
	GetUser(id uuid.UUID) (models.User, error)
	UpdateUser(id uuid.UUID, user models.User) error
	DeleteUser(id uuid.UUID) error
	ListUsers() ([]models.User, error)
	GetUserByEmail(email string) (models.User, error)
}

type userService struct {
	repo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) CreateUser(user models.User) (models.User, error) {
	id, err := s.repo.CreateUser(user)
	if err != nil {
		return models.User{}, err
	}
	return s.repo.GetUser(id)
}

func (s *userService) GetUser(id uuid.UUID) (models.User, error) {
	return s.repo.GetUser(id)
}

func (s *userService) UpdateUser(id uuid.UUID, user models.User) error {
	return s.repo.UpdateUser(id, user)
}

func (s *userService) DeleteUser(id uuid.UUID) error {
	return s.repo.DeleteUser(id)
}

func (s *userService) ListUsers() ([]models.User, error) {
	return s.repo.ListUsers()
}

func (s *userService) GetUserByEmail(email string) (models.User, error) {
	return s.repo.GetUserByEmail(email)
}
