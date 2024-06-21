package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"social-network-otus/internal/session"

	"social-network-otus/internal/auth"
	"social-network-otus/internal/rest/response"
	"social-network-otus/internal/user"
	"social-network-otus/internal/utils"
	"social-network-otus/internal/validator"
)

type UserId struct {
	Id string `uri:"id" json:"id" binding:"required,uuid"`
}

type FindUser struct {
	Name    string `form:"name" json:"name" binding:"required_without=Surname" _required_without:"$field should be filled if surname is empty"`
	Surname string `form:"surname" json:"surname" binding:"required_without=Name" _required_without:"$field should be filled if name is empty"`
}

type ChangePassword struct {
	OldPassword     string `json:"old_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required"`
	ConfirmPassword string `json:"confirm_password" binding:"required"`
}

type UserHandler interface {
	UserPage(c *gin.Context)
	FindUsers(c *gin.Context)
	UpdatePassword(c *gin.Context)
}

type UserHandlerInstance struct {
	authService *auth.AuthService
	userService *user.Service
	session     *session.SessionStorage
	response    *response.ResponseFactory
}

func NewUserHandler(authService *auth.AuthService, session *session.SessionStorage, userService *user.Service, responseFactory *response.ResponseFactory) *UserHandlerInstance {
	return &UserHandlerInstance{authService: authService, session: session, userService: userService, response: responseFactory}
}

func (h *UserHandlerInstance) UserPage(c *gin.Context) {
	var input UserId
	newValidator := validator.NewValidator(input)

	if err := c.ShouldBindUri(&input); err != nil {
		h.response.BadRequest(c, newValidator.DecryptErrors(err).(response.F))
		return
	}

	//https://github.com/gin-gonic/gin/issues/2423
	userUUID, err := uuid.Parse(input.Id)
	if err != nil {
		h.response.BadRequest(c, response.F{"id": "invalid user id"})
	}

	user, appError := h.userService.GetUserById(userUUID.String())

	if appError != nil {
		h.response.FromAppError(c, appError, utils.Ptr("id"))
		return
	}

	h.response.Ok(c, user)
}

func (h *UserHandlerInstance) UpdatePassword(c *gin.Context) {
	var input ChangePassword
	user := h.session.GetAuthenticatedUser()

	validator := validator.NewValidator(input)

	if err := c.ShouldBindJSON(&input); err != nil {
		errors := validator.DecryptErrors(err)
		h.response.BadRequest(c, errors.(response.F))
		return
	}

	if input.NewPassword != input.ConfirmPassword {
		h.response.BadRequest(c, response.F{"confirm_password": "error confirming password"})
		return
	}

	user, appErr := h.userService.UpdatePassword(user, input.OldPassword, input.NewPassword)

	if appErr != nil {
		h.response.FromAppError(c, appErr, utils.Ptr("confirm_password"))
		return
	}

	h.response.OkWithMessage(c, "password was updated")
}

func (h *UserHandlerInstance) FindUsers(c *gin.Context) {
	var input FindUser
	validator := validator.NewValidator(input)

	if err := c.ShouldBindQuery(&input); err != nil {
		h.response.BadRequest(c, validator.DecryptErrors(err).(response.F))
		return
	}

	users, appError := h.userService.FindUsers(input.Name, input.Surname)

	if appError != nil {
		h.response.FromAppError(c, appError, nil)
		return
	}

	h.response.Ok(c, response.F{"users": users})
}
