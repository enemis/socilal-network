package post

import (
	"time"

	"github.com/pkg/errors"

	"social-network-otus/internal/app_error"
	"social-network-otus/internal/database"
)

type PostRepository interface {
	GetPost(postId string) (*Post, *app_error.AppError)
	CreatePost(post *Post) (*Post, *app_error.AppError)
	UpdatePost(post *Post) (*Post, *app_error.AppError)
}

type PostRepositoryInstance struct {
	db *database.DatabaseStack
}

func NewPostRepository(databaseStack *database.DatabaseStack) *PostRepositoryInstance {
	return &PostRepositoryInstance{db: databaseStack}
}

func (r *PostRepositoryInstance) GetPost(postId string) (*Post, *app_error.AppError) {
	var post Post
	err := r.db.GetReadConnection().Get(&post, "SELECT * FROM posts WHERE id=$1 LIMIT 1", postId)

	if err != nil {
		return nil, app_error.NewBadRequestFromError(errors.New("post not found"))
	}

	return &post, nil
}

func (r *PostRepositoryInstance) CreatePost(post *Post) (*Post, *app_error.AppError) {
	query := "INSERT INTO posts (id, user_id, title, post, created_at, update_at, status) VALUES ($1, $2, $3, $4, $5, $6, $7)"

	err := r.db.GetWriteConnection().QueryRow(query, post.Id, post.UserId, post.Title, post.Post, post.CreatedAt, post.UpdatedAt).Scan()
	if err != nil {
		return nil, app_error.NewInternalServerError(errors.Wrap(err, "error create user"))
	}

	post, apperr := r.GetPost(post.Id.String())
	if apperr != nil {
		return nil, apperr
	}
	return post, apperr
}

func (r *PostRepositoryInstance) UpdatePost(post *Post) (*Post, *app_error.AppError) {
	query := "UPDATE posts " +
		"SET (title, post, update_at, status) " +
		"VALUES ($1, $2, $3, $4) WHERE id=$5"

	post.UpdatedAt = time.Now()
	err := r.db.GetWriteConnection().QueryRow(query, post.Title, post.Post, post.UpdatedAt, post.Status, post.Id).Scan()
	if err != nil {
		return nil, app_error.NewInternalServerError(err)
	}

	return post, nil
}
