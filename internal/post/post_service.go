package post

import (
	"fmt"
	"social-network-otus/internal/app_error"
	"social-network-otus/internal/logger"
	"social-network-otus/internal/session"
	"social-network-otus/internal/user"
)

type PostService struct {
	feedService *FeedService
	userService *user.Service
	repository  PostRepository
	session     *session.SessionStorage
	logger      logger.LoggerInterface
}

func NewPostService(repository PostRepository, logger logger.LoggerInterface, session *session.SessionStorage, feedService *FeedService, userService *user.Service) *PostService {
	return &PostService{
		logger:      logger,
		repository:  repository,
		feedService: feedService,
		userService: userService,
		session:     session,
	}
}

func (s *PostService) GetPost(postId string) (*Post, *app_error.AppError) {
	post, appError := s.repository.GetPost(postId)

	if appError != nil {
		return nil, appError
	}

	user := s.session.GetAuthenticatedUser()
	if post.UserId.String() != user.Id.String() && post.Status != Published {
		return nil, app_error.NewForbiddenFromError(fmt.Errorf("access denied to view post"))
	}

	return post, nil
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

		s.feedService.NotifyFriends(usr, post)
	}

	return post, nil
}

func (s *PostService) UpdatePost(post *Post) (*Post, *app_error.AppError) {
	user := s.session.GetAuthenticatedUser()
	if post.UserId.String() != user.Id.String() {
		return nil, app_error.NewBadRequestFromError(fmt.Errorf("current user doenst have permission to edit this post"))
	}
	post, err := s.repository.UpdatePost(post)
	if err != nil {
		return nil, err
	}

	if post.Status == Published {
		usr, appErr := s.userService.GetUserById(post.UserId.String())
		if appErr != nil {
			s.logger.Error("error fetch user", appErr.OriginalError(), nil)
			return post, nil
		}

		s.feedService.NotifyFriends(usr, post)
	}

	return post, nil
}

func (s *PostService) DeletePost(postId string) *app_error.AppError {
	post, err := s.repository.GetPost(postId)
	if err != nil {
		return err
	}
	user := s.session.GetAuthenticatedUser()

	if post.UserId.String() != user.Id.String() {
		return app_error.NewForbiddenFromError(fmt.Errorf("current user doenst have permission to delete post"))
	}

	appErr := s.repository.DeletePost(postId)
	if appErr != nil {
		return appErr
	}

	if post.Status == Published {
		usr, appErr := s.userService.GetUserById(post.UserId.String())
		if appErr != nil {
			s.logger.Error("error fetch user", appErr.OriginalError(), nil)
			return nil
		}

		s.feedService.NotifyFriends(usr, nil)
	}

	return nil
}
