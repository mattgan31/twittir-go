package services

import (
	"errors"

	"twittir-go/internal/domain"
	"twittir-go/internal/helpers"
	"twittir-go/internal/repositories"
)

type UserService interface {
	Register(username, email, fullname, password, passwordVerif string) (*domain.User, error)
	Login(username, password string) (string, error)
	GetUserByID(id int) (*domain.User, error)
	SearchUserByUsername(username string) ([]domain.User, error)
	UpdateProfile(userID int, updateUser *domain.User) (*domain.User, error)
}

type userService struct {
	userRepo repositories.UserRepository
}

func NewUserService(userRepo repositories.UserRepository) UserService {
	return &userService{userRepo}
}

func (s *userService) Register(username, email, fullname, password, passwordVerif string) (*domain.User, error) {
	// Check if password and password verification are the same
	if password != passwordVerif {
		return nil, errors.New("password and password verification must be the same")
	}

	// Hash password
	hashedPassword := helpers.HashPass(password)

	// Create user
	user := domain.User{
		Username: username,
		Email:    email,
		FullName: fullname,
		Password: hashedPassword,
	}

	// Save user to database
	err := s.userRepo.SaveNewUser(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *userService) Login(username, password string) (string, error) {
	user, err := s.userRepo.FindByUsername(username)
	if err != nil {
		return "", err
	}

	// Compare passwords
	if !helpers.ComparePass([]byte(user.Password), []byte(password)) {
		return "", errors.New("invalid username/password")
	}

	// Generate token
	token := helpers.GenerateToken(user.ID, user.Username)
	return token, nil
}

func (s *userService) GetUserByID(id int) (*domain.User, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Return user
	return user, nil
}

func (s *userService) SearchUserByUsername(username string) ([]domain.User, error) {
	users, err := s.userRepo.SearchByUsername(username)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (s *userService) UpdateProfile(userID int, updateUser *domain.User) (*domain.User, error) {

	user, err := s.userRepo.UpdateProfile(userID, updateUser)
	if err != nil {
		return nil, err
	}

	return user, nil
}
