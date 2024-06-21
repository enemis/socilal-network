package post

import (
	"social-network-otus/internal/user"

	"social-network-otus/internal/app_error"
	"social-network-otus/internal/database"
)

type FeedRepository interface {
	FeedPosts(user *user.User, limit, offset uint) ([]*Post, *app_error.AppError)
}

type FeedRepositoryInstance struct {
	db *database.DatabaseStack
}

func NewFeedRepository(databaseStack *database.DatabaseStack) *FeedRepositoryInstance {
	return &FeedRepositoryInstance{db: databaseStack}
}

// todo on production replace to cursors instead
func (r *FeedRepositoryInstance) FeedPosts(user *user.User, limit, offset uint) ([]*Post, *app_error.AppError) {
	var results = make([]*Post, 0, limit)

	query := "WITH user_friends AS (" +
		"SELECT friend_id FROM friends WHERE user_id = $1" +
		")" +
		"SELECT * FROM posts p " +
		"WHERE deleted_at IS NULL AND p.status = 'published' " +
		"AND user_id IN (select friend_id FROM user_friends) " +
		"ORDER BY p.updated_at DESC LIMIT $2 OFFSET $3"

	rows, err := r.db.GetReadConnection().Queryx(query, user.Id, limit, offset)

	if err != nil {
		return nil, app_error.NewInternalServerError(err)
	}

	for rows.Next() {
		var post Post
		err := rows.StructScan(&post)
		if err != nil {
			return nil, app_error.NewInternalServerError(err)
		}

		results = append(results, &post)
	}

	return results, nil
}
