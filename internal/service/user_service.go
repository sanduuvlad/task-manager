package service

import (
	"errors"
	"myprojects/internal/models"
	"myprojects/internal/repository"
)

var ErrUsersNotFounf = errors.New("users not found")
var ErrInvalidTask = errors.New("Invalid task")

type UserService struct {
	Repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{Repo: repo}
}

func (s *UserService) GetUserByID(id int64) (models.User, error) {
	return s.Repo.GetUserByID(id)
}

func (s *UserService) CreateUser(user models.User) error {
	return nil
}

func (s *UserService) GetAllUsers() ([]models.User, error) {
	users, err := s.Repo.GetAllUsers()
	if err != nil {
		return nil, err
	}

	if len(users) == 0 {
		return nil, ErrUsersNotFounf
	}

	return users, nil
}
