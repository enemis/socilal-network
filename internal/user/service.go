package user

import (
	"github.com/google/uuid"

	"social-network-otus/internal/app_error"
	"social-network-otus/internal/token"
)

type Service struct {
	repository        UserRepository
	passwordGenerator token.PasswordGenerator
}

func NewUserService(repository UserRepository, passwordGenerator token.PasswordGenerator) *Service {
	return &Service{
		repository:        repository,
		passwordGenerator: passwordGenerator,
	}
}

func (s *Service) GetUserByEmail(email string) (*User, *app_error.AppError) {
	return s.repository.GetUserByEmail(email)
}

func (s *Service) GetUserById(userId string) (*User, *app_error.AppError) {
	return s.repository.GetUserById(userId)
}

func (s *Service) CreateUser(user *User) (*uuid.UUID, *app_error.AppError) {
	userId := uuid.New()
	pass, err := s.passwordGenerator.PasswordHash(userId, user.Password)

	if err != nil {
		return nil, app_error.NewInternalServerError(err)
	}

	user.Password = pass
	return s.repository.CreateUser(user)
}

func (s *Service) FindUsers(name, surname string) ([]*User, *app_error.AppError) {
	return s.repository.FindUsers(name, surname)
}
