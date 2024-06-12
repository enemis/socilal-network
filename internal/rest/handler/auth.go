package handler

import (
	"github.com/gin-gonic/gin"

	"social-network-otus/internal/auth"
	"social-network-otus/internal/rest/response"
	"social-network-otus/internal/user"
	"social-network-otus/internal/validator"
)

type SignInInput struct {
	Email    string `form:"email" binding:"required,email"`
	Password string `form:"password" binding:"required"`
}

type AuthHandler interface {
	SignIn(c *gin.Context)
	SignUp(c *gin.Context)
}

type AuthHandlerInstance struct {
	authService *auth.AuthService
	userService *user.Service
	response    *response.ResponseFactory
}

func NewAuthHandler(authService *auth.AuthService, userService *user.Service, responseFactory *response.ResponseFactory) *AuthHandlerInstance {
	return &AuthHandlerInstance{authService: authService, userService: userService, response: responseFactory}
}

func (h *AuthHandlerInstance) SignIn(c *gin.Context) {
	var input SignInInput
	newValidator := validator.NewValidator(input)
	if err := c.ShouldBind(&input); err != nil {
		h.response.BadRequest(c, response.F{"errors": newValidator.DecryptErrors(err)})
	}

	token, err := h.authService.Login(c.Request.Context(), input.Email, input.Password)

	if err != nil {
		h.response.BadRequest(c, response.F{"email": "invalid email"})
		return
	}

	h.response.Ok(c, response.F{"token": token})
}

func (h *AuthHandlerInstance) SignUp(c *gin.Context) {
	var input user.User
	validator := validator.NewValidator(input)

	if err := c.ShouldBindJSON(&input); err != nil {
		h.response.BadRequest(c, response.F{"errors": validator.DecryptErrors(err)})
		return
	}

	uuid, err := h.userService.CreateUser(&input)

	if err != nil {
		h.response.InternalServerError(c, err)
		return
	}

	h.response.Created(c, uuid)
}
