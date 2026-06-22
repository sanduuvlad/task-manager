package service

import (
	"errors"
	"myprojects/internal/models"
)

var ErrUsersNotFound = errors.New("users not found")
var ErrInvalidTask = errors.New("Invalid task")
var ErrUserAlreadyExists = errors.New("not exist email")

type UserService struct {
	Repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{
		Repo: repo,
	}
}

type UserRepository interface {
	GetUserByID(id int64) (models.User, error)
	CreateUser(user models.User) error
	GetAllUsers() ([]models.User, error)
}

func (s *UserService) GetUserByID(id int64) (models.User, error) {
	return s.Repo.GetUserByID(id)
}

func (s *UserService) GetAllUsers() ([]models.User, error) {
	users, err := s.Repo.GetAllUsers()
	if err != nil {
		return nil, err
	}

	if len(users) == 0 {
		return nil, ErrUsersNotFound
	}

	return users, nil
}
