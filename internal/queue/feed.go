package queue

import (
	"social-network-otus/internal/config"
	"time"
)

func initFeedQueueClient(factory *ClientFactory, config *config.Config) *Client {
	return factory.NewClient(Queue{
		Name:      config.QueueWarmUpFeed,
		AutoAck:   true,
		Durable:   false,
		DelUnused: false,
		Exclusive: false,
		NoWait:    false,
		Arg:       nil,
	})
}

type FeedWarmUpQueueItem struct {
	UserId string    `json:"user_id"`
	Date   time.Time `json:"date"`
}
