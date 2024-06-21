package post

import (
	"context"
	"fmt"
	"github.com/goccy/go-json"
	strip "github.com/grokify/html-strip-tags-go"
	"go.uber.org/fx"
	"social-network-otus/internal/app_error"
	"social-network-otus/internal/config"
	"social-network-otus/internal/friend"
	"social-network-otus/internal/logger"
	"social-network-otus/internal/queue"
	"social-network-otus/internal/user"
	"social-network-otus/internal/utils"
	"time"
)

type FeedService struct {
	cache            *Cache
	repository       FeedRepository
	friendRepository friend.FriendRepository
	config           *config.Config
	logger           logger.LoggerInterface
	userService      *user.Service
	FeedWarmupQueue  *queue.Client
}

func NewFeedService(cache *Cache, userService *user.Service, FeedWarmupQueue *queue.Client, config *config.Config, logger logger.LoggerInterface, repository FeedRepository, friendRepository friend.FriendRepository) *FeedService {
	return &FeedService{
		cache:            cache,
		config:           config,
		friendRepository: friendRepository,
		repository:       repository,
		logger:           logger,
		userService:      userService,
		FeedWarmupQueue:  FeedWarmupQueue,
	}
}

func (s *FeedService) Feed(user *user.User, limit, offset uint) (*Feed, *app_error.AppError) {
	for {
		feed, err := s.cache.GetFeed(user.Id.String())
		if err == nil && feed.Status == Ready {
			return feed, nil
		}

		if feed.Status == Fetching {
			<-time.After(time.Second)
			continue
		}
		feed = &Feed{Status: Fetching}
		s.cache.SetFeed(user.Id.String(), feed, 5*time.Second)
		posts, appError := s.repository.FeedPosts(user, limit, offset)
		if appError != nil {
			return nil, appError
		}
		feed = &Feed{Status: Ready}
		for _, post := range posts {
			feed.Items = append(feed.Items, postToFeedItem(post))
		}

		s.cache.SetFeed(user.Id.String(), feed, s.config.FeedCache)

		return s.limitFees(feed, limit, offset), nil
	}
}

func (s *FeedService) limitFees(feed *Feed, limit, offset uint) *Feed {
	lenSlice := uint(len(feed.Items))

	from := limit * offset
	to := from + limit
	if from > lenSlice {
		feed.Items = make([]*FeedItem, 0)
		return feed
	}
	if to > lenSlice {
		to = lenSlice
	}

	feed.Items = feed.Items[from:to]
	return feed
}

func (s *FeedService) ScheduleFriendsFeedWarming(user *user.User) *app_error.AppError {
	friends, err := s.friendRepository.GetFriends(user)
	if err != nil {
		return err
	}

	for _, friend := range friends {
		data, err := json.Marshal(queue.FeedWarmUpQueueItem{
			UserId: friend.Id.String(),
			Date:   time.Now(),
		})
		if err != nil {
			s.logger.Error("Marshaling error", err, nil)
		}
		s.FeedWarmupQueue.Push(data)
	}

	return nil
}

func (s *FeedService) updateFeed(userId string) {
	s.cache.Clear(userId)
	usr, err := s.userService.GetUserById(userId)

	if err != nil {
		s.logger.Error(err.Error(), err.OriginalError(), nil)
	}

	sevenDaysAgo := time.Now().AddDate(0, 0, -7)
	if usr.LastOnline.Unix() > sevenDaysAgo.Unix() {
		s.Feed(usr, s.config.FeedLastCount, 0)
	}
}

func InitHooks(lc fx.Lifecycle, feedWarmUpConsumer *FeedWarmupConsumer) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				feedWarmUpConsumer.Consume()
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return nil
		},
	})
}

func postToFeedItem(post *Post) *FeedItem {
	return &FeedItem{
		Title:  post.Title,
		Path:   fmt.Sprintf("posts/%s", post.Id.String()),
		Teaser: utils.EllipticalTruncate(strip.StripTags(post.Post), 150),
	}
}
