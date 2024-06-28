package notifier_ws

import (
	"github.com/gorilla/websocket"
	"social-network-otus/internal/logger"
	"social-network-otus/internal/user"
)

const New_post_message = iota + 1

type NotifyMessage interface {
	GetMessage() (*string, error)
	GetType() int
}

type Notifier struct {
	users  map[string]*websocket.Conn
	logger logger.LoggerInterface
}

func NewNotifier(logger logger.LoggerInterface) *Notifier {
	return &Notifier{users: make(map[string]*websocket.Conn), logger: logger}
}

func (n *Notifier) Subscribe(user *user.User, conn *websocket.Conn) {
	n.users[user.Id.String()] = conn
}

func (n *Notifier) IsUserConnected(user *user.User) bool {
	if n.users[user.Id.String()] == nil {
		return false
	}

	return true
}

func (n *Notifier) UnSubscribe(user *user.User) {
	n.users[user.Id.String()] = nil
}

func (n *Notifier) Notify(user *user.User, msg NotifyMessage) (bool, error) {
	if n.users[user.Id.String()] == nil {
		return false, nil
	}

	msgString, err := msg.GetMessage()
	if err != nil {
		n.logger.Error(err.Error(), err, nil)
		return false, err
	}

	err = n.users[user.Id.String()].WriteMessage(msg.GetType(), []byte(*msgString))
	if err != nil {
		n.logger.Error(err.Error(), err, nil)
		return false, err
	}

	return true, nil
}
