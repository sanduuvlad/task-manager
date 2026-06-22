package service

import (
	"errors"
	"myprojects/internal/dto"
	"myprojects/internal/models"
	"myprojects/internal/repository"

	"golang.org/x/crypto/bcrypt"
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
	GetUserByEmail(email string) (models.User, error)
}

func (s *UserService) GetUserByID(id int64) (models.User, error) {
	return s.Repo.GetUserByID(id)
}

func (s *UserService) CreateUser(req dto.CreateUserRequest) error {
	if req.Username == "" {
		return errors.New("username is required")
	}

	if req.Email == "" {
		return errors.New("email is required")
	}

	if len(req.Password) < 8 {
		return errors.New("password must be at least 8 characters")
	}

	_, err := s.Repo.GetUserByEmail(req.Email)
	if err != nil {
		return ErrUserAlreadyExists
	}

	if !errors.Is(err, repository.ErrUserNotFound) {
		return err
	}

	hash, err := bcrypt.GenerateFromPassword(
		[]byte(req.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return err
	}

	user := models.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hash),
	}

	return s.Repo.CreateUser(user)
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
