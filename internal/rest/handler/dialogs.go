package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"social-network-otus/internal/client"
	pb "social-network-otus/internal/client/proto"
	"social-network-otus/internal/rest/response"
	"social-network-otus/internal/session"
	"social-network-otus/internal/user"
	"social-network-otus/internal/validator"
)

type DialogsHandler interface {
	SendMessage(c *gin.Context)
	GetDialog(c *gin.Context)
}

type MessageInput struct {
	Recipient uuid.UUID `json:"recipient" binding:"required"`
	Message   string    `json:"message" binding:"required"`
}

type DialogsHandlerInstance struct {
	session     *session.SessionStorage
	userService *user.Service
	response    *response.ResponseFactory
	client      *client.Client
}

func NewDialogsHandler(userService *user.Service, client *client.Client, responseFactory *response.ResponseFactory, session *session.SessionStorage) *DialogsHandlerInstance {
	return &DialogsHandlerInstance{session: session, client: client, userService: userService, response: responseFactory}
}

func (d *DialogsHandlerInstance) SendMessage(c *gin.Context) {
	var input MessageInput

	validator := validator.NewValidator(input)

	if err := c.ShouldBind(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, validator.DecryptErrors(err).(map[string]interface{}))
		return
	}

	response, err := d.client.ClientGRPC.SendMessage(context.Background(), &pb.SendMessageDialogRequest{
		Sender:   &pb.User{Id: &pb.UUID{Value: d.session.GetAuthenticatedUser().Id.String()}},
		Reciever: &pb.User{Id: &pb.UUID{Value: input.Recipient.String()}},
		Message:  input.Message,
	})

	if err != nil {
		d.response.BadRequest(c, map[string]interface{}{"error": err.Error()})
	}

	d.response.Ok(c, map[string]interface{}{"message": response.Message})
}

func (d *DialogsHandlerInstance) GetDialog(c *gin.Context) {
	//friend, err := h.resolveFriendFromRequest(c)
	//if err != nil {
	//	return
	//}
	//
	//user := h.session.GetAuthenticatedUser()
	//
	//if err != nil {
	//	return
	//}
	//
	//err = h.friendService.RemoveFriend(user, friend)
	//
	//h.response.OkWithMessage(c, fmt.Sprintf("User %s %s has been removed from your friends list", friend.Surname, friend.Name))
}
