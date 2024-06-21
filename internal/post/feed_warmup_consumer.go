package post

import (
	"fmt"
	"github.com/goccy/go-json"
	amqp "github.com/rabbitmq/amqp091-go"
	"social-network-otus/internal/logger"
	"social-network-otus/internal/queue"
	"time"
)

type FeedWarmupConsumer struct {
	logger          logger.LoggerInterface
	feedWarmupQueue *queue.Client
	feedService     *FeedService
}

func NewFeedWarmupConsumer(logger logger.LoggerInterface, feedService *FeedService, FeedWarmupQueue *queue.Client) *FeedWarmupConsumer {
	return &FeedWarmupConsumer{
		logger:          logger,
		feedWarmupQueue: FeedWarmupQueue,
		feedService:     feedService,
	}
}

func (s *FeedWarmupConsumer) Consume() {
	//wait for establishing connection
	<-time.After(time.Second * 2)
	deliveries, err := s.feedWarmupQueue.Consume()
	if err != nil {
		s.logger.Fatal("could not start consuming: %s\n", err, nil)
		return
	}

	// This channel will receive a notification when a channel closed event
	// happens. This must be different from Client.notifyChanClose because the
	// library sends only one notification and Client.notifyChanClose already has
	// a receiver in handleReconnect().
	// Recommended to make it buffered to avoid deadlocks
	chClosedCh := make(chan *amqp.Error, 1)
	s.feedWarmupQueue.Channel.NotifyClose(chClosedCh)

	for {
		select {
		case amqErr := <-chClosedCh:
			// This case handles the event of closed channel e.g. abnormal shutdown
			s.logger.Warn(fmt.Sprintf("AMQP Channel closed due to: %s\n", amqErr), nil)

			deliveries, err = s.feedWarmupQueue.Consume()
			if err != nil {
				// If the AMQP channel is not ready, it will continue the loop. Next
				// iteration will enter this case because chClosedCh is closed by the
				// library
				s.logger.Warn(fmt.Sprintf("error trying to consume, will try again"), nil)
				continue
			}

			// Re-set channel to receive notifications
			// The library closes this channel after abnormal shutdown
			chClosedCh = make(chan *amqp.Error, 1)
			s.feedWarmupQueue.Channel.NotifyClose(chClosedCh)

		case delivery := <-deliveries:
			var queueItem queue.FeedWarmUpQueueItem
			err = json.Unmarshal(delivery.Body, &queueItem)
			if err != nil {
				delivery.Nack(false, false)
				s.logger.Error(fmt.Sprintf("error unmarshal feed warmup queue: %s", delivery.Body), err, nil)
			}
			s.feedService.updateFeed(queueItem.UserId)

		}
	}
}
