package queue

import (
	"context"
	"fmt"
	"go.uber.org/fx"
	"social-network-otus/internal/config"
	"social-network-otus/internal/logger"
)

var Module = fx.Options(
	fx.Provide(
		NewFactory,
		initClients,
	),
	fx.Invoke(
		InitHooks,
	),
)

type QueueCLients struct {
	fx.Out
	FeedWarmupQueue *Client
}

func initClients(factory *ClientFactory, config *config.Config) QueueCLients {
	client := initFeedQueueClient(factory, config)

	return QueueCLients{
		FeedWarmupQueue: client,
	}
}

func InitHooks(lc fx.Lifecycle, logger logger.LoggerInterface, FeedWarmupQueue *Client) {
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			logger.Info(fmt.Sprintf("Stop %s queue client", FeedWarmupQueue.Queue.Name), nil)
			err := FeedWarmupQueue.Close()
			if err != nil {
				logger.Error(fmt.Sprintf("Stop %s queue client error", FeedWarmupQueue.Queue.Name), err, nil)
			}
			return err
		},
	})
}
