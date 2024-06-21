package post

import (
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
	Title  string
	Path   string
	Teaser string
}

type Feed struct {
	Items  []*FeedItem
	Status StatusFeed
}
