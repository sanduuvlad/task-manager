package service

import (
	"myprojects/internal/models"
	"myprojects/internal/repository"
)

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
	return s.Repo.GetAllUsers()
}
