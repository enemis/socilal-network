package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"social-network-otus/internal/logger"
	"social-network-otus/internal/notifier_ws"
	"social-network-otus/internal/rest/response"
	"social-network-otus/internal/session"
)

type WSHandler interface {
	Connect(c *gin.Context)
}

type WSHandlerInstance struct {
	logger          logger.LoggerInterface
	responseFactory *response.ResponseFactory
	session         *session.SessionStorage
	upgrader        websocket.Upgrader
	notifier        *notifier_ws.Notifier
}

func NewWSHandlerInstance(responseFactory *response.ResponseFactory, session *session.SessionStorage, logger logger.LoggerInterface, notifier *notifier_ws.Notifier) *WSHandlerInstance {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	return &WSHandlerInstance{logger: logger, session: session, responseFactory: responseFactory, upgrader: upgrader, notifier: notifier}
}

func (h *WSHandlerInstance) Connect(c *gin.Context) {
	conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		h.responseFactory.InternalServerError(c, fmt.Errorf("failed to upgrade connection to WebSocket"))
		return
	}
	user := h.session.GetAuthenticatedUser()
	fmt.Println("Client user id=%s connected to WebSocket", user.Id.String())

	h.notifier.Subscribe(user, conn)

	defer func() {
		conn.Close()
		h.notifier.UnSubscribe(user)
	}()

	// Infinite loop to handle incoming messages
	for {
		// Read message from the WebSocket connection
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Error reading message:", err)
			return
		}

		// Print the received message to the console
		fmt.Printf("Received message: %s\n", p)

		// Broadcast the received message to all connected clients
		err = conn.WriteMessage(messageType, p)
		if err != nil {
			fmt.Println("Error broadcasting message:", err)
			return
		}
	}
}
