package user

import (
	"fmt"
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
	user, err := s.repository.GetUserById(userId)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Service) CreateUser(user *User) (*uuid.UUID, *app_error.AppError) {
	userId := uuid.New()
	pass, err := s.passwordGenerator.PasswordHash(userId, user.Password)

	if err != nil {
		return nil, app_error.NewInternalServerError(err)
	}

	user.Password = pass
	user.Id = userId
	return s.repository.CreateUser(user)
}

func (s *Service) UpdatePassword(user *User, oldPassword, newPassword string) (*User, *app_error.AppError) {
	//move to custom validation rule
	if !s.passwordGenerator.CompareHashAndPassword(user.Id, user.Password, oldPassword) {
		return nil, app_error.NewBadRequestFromError(fmt.Errorf("invalid password"))
	}

	pass, err := s.passwordGenerator.PasswordHash(user.Id, newPassword)

	if err != nil {
		return nil, app_error.NewInternalServerError(err)
	}

	user.Password = pass
	apperr := s.repository.UpdateUser(user)

	if apperr != nil {
		return nil, apperr
	}

	return user, nil
}

func (s *Service) FindUsers(name, surname string) ([]*User, *app_error.AppError) {
	return s.repository.FindUsers(name, surname)
}
