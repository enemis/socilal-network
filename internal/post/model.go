package post

import (
	"github.com/goccy/go-json"
	"social-network-otus/internal/notifier_ws"
	"social-network-otus/internal/utils"
	"time"

	"github.com/google/uuid"
)

type PostStatus string

const (
	Draft     PostStatus = "draft"
	Published PostStatus = "published"
)

type Post struct {
	Id        uuid.UUID  `json:"id" db:"id"`
	UserId    uuid.UUID  `json:"user_id" binding:"required" db:"user_id"`
	Title     string     `json:"title" binding:"required,alphanum"`
	Post      string     `json:"post" binding:"required,alphanum"`
	Status    PostStatus `json:"status" binding:"required"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
}

type StatusFeed uint8

const (
	NotFound StatusFeed = iota
	Fetching
	Ready
)

type FeedItem struct {
	Title  string `json:"title"`
	Path   string `json:"path"`
	Teaser string `json:"teaser"`
}

type Feed struct {
	Items  []*FeedItem
	Status StatusFeed
}

type PostWSMessage struct {
	FeedItem
	UserName string `json:"user_name"`
}

func (s PostWSMessage) GetMessage() (*string, error) {
	bytes, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}

	return utils.Ptr(string(bytes)), nil
}

func (s PostWSMessage) GetType() int {
	return notifier_ws.New_post_message
}
