package post

import (
	"github.com/gin-gonic/gin"

	"social-network-otus/internal/app_error"
)

type PostService struct {
	repository PostRepository
}

func NewPostService(repository PostRepository) *PostService {
	return &PostService{
		repository: repository,
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

	return post, nil
}

func (s *PostService) UpdatePost(post *Post) (*Post, *app_error.AppError) {
	post, err := s.repository.UpdatePost(post)
	if err != nil {
		return nil, err
	}
	return post, nil
}

func (s *PostService) DeletePost(c *gin.Context) {
	post, err := s.repository.UpdatePost(post)
}
