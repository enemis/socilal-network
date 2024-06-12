package friend

import (
	"social-network-otus/internal/app_error"
	"social-network-otus/internal/user"
)

type FriendService struct {
	repository FriendRepository
}

func NewFriendService(repository FriendRepository) *FriendService {
	return &FriendService{
		repository: repository,
	}
}

func (s *FriendService) AddFriend(user *user.User, friend *user.User) *app_error.AppError {

	return s.repository.AddFriend(user, friend)
}

func (s *FriendService) RemoveFriend(user *user.User, friend *user.User) *app_error.AppError {

	return s.repository.RemoveFriend(user, friend)
}
