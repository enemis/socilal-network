package handler

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"social-network-otus/internal/friend"
	"social-network-otus/internal/rest/response"
	"social-network-otus/internal/session"
	"social-network-otus/internal/user"
	"social-network-otus/internal/validator"
)

type FriendHandler interface {
	AddFriend(c *gin.Context)
	RemoveFriend(c *gin.Context)
}

type FriendId struct {
	Id string `form:"friend_id" json:"friend_id" binding:"required,uuid"`
}

type FriendHandlerInstance struct {
	userService   *user.Service
	friendService *friend.FriendService
	session       *session.SessionStorage
	response      *response.ResponseFactory
}

func NewFriendHandler(userService *user.Service, friendService *friend.FriendService, responseFactory *response.ResponseFactory, session *session.SessionStorage) *FriendHandlerInstance {
	return &FriendHandlerInstance{session: session, friendService: friendService, userService: userService, response: responseFactory}
}

func (h *FriendHandlerInstance) AddFriend(c *gin.Context) {
	friend, err := h.resolveFriendFromRequest(c)
	if err != nil {
		return
	}

	user := h.session.GetAuthenticatedUser()

	httpError := h.friendService.AddFriend(user, friend)
	if httpError != nil {
		h.response.FromAppError(c, httpError, nil)
		return
	}

	h.response.OkWithMessage(c, fmt.Sprintf("User %s %s now is in your friends list", friend.Surname, friend.Name))
}

func (h *FriendHandlerInstance) RemoveFriend(c *gin.Context) {
	friend, err := h.resolveFriendFromRequest(c)
	if err != nil {
		return
	}

	user := h.session.GetAuthenticatedUser()

	if err != nil {
		return
	}

	err = h.friendService.RemoveFriend(user, friend)

	h.response.OkWithMessage(c, fmt.Sprintf("User %s %s has been removed from your friends list", friend.Surname, friend.Name))
}

func (h *FriendHandlerInstance) resolveFriendFromRequest(c *gin.Context) (*user.User, error) {
	var input FriendId
	validator := validator.NewValidator(input)

	if err := c.ShouldBind(&input); err != nil {
		h.response.BadRequest(c, validator.DecryptErrors(err).(response.F))
		return nil, err
	}

	userUUID, err := uuid.Parse(input.Id)
	if err != nil {
		h.response.BadRequest(c, response.F{"friend_id": "friend not found"})
	}

	user, httpError := h.userService.GetUserById(userUUID.String())

	if httpError != nil {
		h.response.BadRequest(c, response.F{"friend_id": "friend not found"})
		return nil, err
	}

	return user, nil
}
