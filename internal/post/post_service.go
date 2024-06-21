package post

import (
	"social-network-otus/internal/app_error"
	"social-network-otus/internal/logger"
	"social-network-otus/internal/user"
)

type PostService struct {
	feedService *FeedService
	userService *user.Service
	repository  PostRepository
	logger      logger.LoggerInterface
}

func NewPostService(repository PostRepository, logger logger.LoggerInterface, feedService *FeedService, userService *user.Service) *PostService {
	return &PostService{
		logger:      logger,
		repository:  repository,
		feedService: feedService,
		userService: userService,
	}
}

func (s *PostService) GetPost(postId string) (*Post, *app_error.AppError) {
	return s.repository.GetPost(postId)
}

func (s *PostService) CreatePost(post *Post) (*Post, *app_error.AppError) {
	post, err := s.repository.CreatePost(post)

	if err != nil {
		return nil, app_error.NewInternalServerError(err)
	}

	if post.Status == Published {
		usr, appErr := s.userService.GetUserById(post.UserId.String())
		if appErr != nil {
			s.logger.Error("error fetch user", appErr.OriginalError(), nil)
			return post, nil
		}

		s.feedService.ScheduleFriendsFeedWarming(usr)
	}

	return post, nil
}

func (s *PostService) UpdatePost(post *Post) (*Post, *app_error.AppError) {
	post, err := s.repository.UpdatePost(post)
	if err != nil {
		return nil, err
	}
	return post, nil
}

func (s *PostService) DeletePost(postId string) *app_error.AppError {
	return s.repository.DeletePost(postId)
}
