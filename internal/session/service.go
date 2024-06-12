package session

import (
	"github.com/google/uuid"

	"social-network-otus/internal/user"
)

type SessionStorage struct {
	user *user.User
}

func NewSessionStorage() *SessionStorage {
	var user *user.User
	return &SessionStorage{user: user}
}

func (s *SessionStorage) IsAuthenticated() bool {
	return s.user.Id != uuid.Nil
}

func (s *SessionStorage) SetAuthenticatedUser(user *user.User) {
	s.user = user
}

func (s *SessionStorage) GetAuthenticatedUser() *user.User {
	return s.user
}
