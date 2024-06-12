package handler

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gitlab.com/pay2play/core-api/pkg/helpers"

	"social-network-otus/internal/post"
	"social-network-otus/internal/rest/response"
	"social-network-otus/internal/session"
	"social-network-otus/internal/validator"
)

type PostHandler interface {
	GetPost(c *gin.Context)
	CreatePost(c *gin.Context)
	UpdatePost(c *gin.Context)
	DeletePost(c *gin.Context)
}

type PostInput struct {
	Title  string          `json:"name" binding:"required" required:"$field is required"`
	Post   string          `json:"post" binding:"required" required:"$field is required"`
	Status post.PostStatus `json:"status" binding:"required,post_status"`
}

type PostHandlerInstance struct {
	response    *response.ResponseFactory
	session     *session.SessionStorage
	postService *post.PostService
}

func NewPostHandler(responseFactory *response.ResponseFactory, session *session.SessionStorage, postService *post.PostService) *PostHandlerInstance {
	return &PostHandlerInstance{session: session, response: responseFactory, postService: postService}
}

func (h *PostHandlerInstance) GetPost(c *gin.Context) {

}

func (h *PostHandlerInstance) CreatePost(c *gin.Context) {
	var input PostInput
	validator := validator.NewValidator(input)

	if err := c.ShouldBind(&input); err != nil {
		h.response.BadRequest(c, validator.DecryptErrors(err).(response.F))
		return
	}

	user := h.session.GetAuthenticatedUser()

	post, apperror := h.postService.CreatePost(&post.Post{
		Id:        uuid.UUID{},
		UserId:    user.Id,
		Title:     input.Title,
		Post:      input.Post,
		Status:    input.Status,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		DeletedAt: nil,
	})

	if apperror != nil {
		h.response.FromAppError(c, apperror, helpers.Ptr("post"))
		return
	}

	h.response.Ok(c, response.F{"post": post})
}

func (h *PostHandlerInstance) UpdatePost(c *gin.Context) {

}

func (h *PostHandlerInstance) DeletePost(c *gin.Context) {

}
