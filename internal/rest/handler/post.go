package handler

import (
	"social-network-otus/internal/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

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
	Status post.PostStatus `json:"status" binding:"required,oneof=draft published" oneof:"$field must be one of draft, published"`
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
	id := c.Param("id")
	post, apperror := h.postService.GetPost(id)
	if apperror != nil {
		h.response.FromAppError(c, apperror, utils.Ptr("id"))
		return
	}

	h.response.Ok(c, response.F{"post": post})
}

func (h *PostHandlerInstance) CreatePost(c *gin.Context) {
	var input PostInput

	validator := validator.NewValidator(input, post.NewPostStatusValidationRule())

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
		h.response.FromAppError(c, apperror, utils.Ptr("post"))
		return
	}

	h.response.Ok(c, response.F{"post": post})
}

func (h *PostHandlerInstance) UpdatePost(c *gin.Context) {
	var input PostInput
	validator := validator.NewValidator(input)

	id := c.Param("id")
	post, apperr := h.postService.GetPost(id)

	if apperr != nil {
		h.response.FromAppError(c, apperr, utils.Ptr("id"))
		return
	}

	if err := c.ShouldBind(&input); err != nil {
		h.response.BadRequest(c, validator.DecryptErrors(err).(response.F))
		return
	}

	post.Status = input.Status
	post.Title = input.Title
	post.Post = input.Post

	post, apperror := h.postService.UpdatePost(post)

	if apperror != nil {
		h.response.FromAppError(c, apperror, utils.Ptr("post"))
		return
	}

	h.response.Ok(c, response.F{"post": post})
}

func (h *PostHandlerInstance) DeletePost(c *gin.Context) {
	id := c.Param("id")
	apperror := h.postService.DeletePost(id)
	if apperror != nil {
		h.response.FromAppError(c, apperror, utils.Ptr("id"))
		return
	}

	h.response.OkWithMessage(c, "deleted")
}
