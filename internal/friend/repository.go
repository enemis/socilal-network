package friend

import (
	"errors"
	"fmt"
	"time"

	"social-network-otus/internal/app_error"
	"social-network-otus/internal/database"
	"social-network-otus/internal/user"
)

type FriendRepository interface {
	GetFriends(user *user.User) ([]*user.User, *app_error.AppError)
	AddFriend(user *user.User, friend *user.User) *app_error.AppError
	RemoveFriend(user *user.User, friend *user.User) *app_error.AppError
}

type FriendRepositoryInstance struct {
	db *database.DatabaseStack
}

func NewFriendRepository(db *database.DatabaseStack) *FriendRepositoryInstance {
	return &FriendRepositoryInstance{db: db}
}

func (r *FriendRepositoryInstance) GetFriends(usr *user.User) ([]*user.User, *app_error.AppError) {
	results := make([]*user.User, 0)

	rows, err := r.db.GetReadConnection().Queryx(
		"SELECT u.* FROM friends f "+
			"INNER JOIN users u ON f.friend_id = u.id "+
			"WHERE user_id=$1", usr.Id)

	if err != nil {
		return nil, app_error.NewInternalServerError(err)
	}

	for rows.Next() {
		var friend user.User
		err := rows.StructScan(&friend)
		if err != nil {
			return nil, app_error.NewInternalServerError(err)
		}

		results = append(results, &friend)
	}

	return results, nil
}

func (r *FriendRepositoryInstance) AddFriend(user *user.User, friend *user.User) *app_error.AppError {
	rows, err := r.db.Slave().Query("SELECT EXISTS(SELECT 1 FROM friends WHERE user_id=$1 AND friend_id=$2)", user.Id, friend.Id)
	if err != nil {
		return app_error.NewInternalServerError(err)
	}

	var exists bool

	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&exists); err != nil {
			return app_error.NewInternalServerError(err)
		}

		if exists {
			return app_error.NewBadRequestFromError(errors.New(fmt.Sprintf("User %s %s already added as friend for user %s %s", friend.Surname, friend.Name, user.Surname, friend.Name, "friend_id")))
		}
	}

	query := "INSERT INTO friends (user_id, friend_id, created_at) VALUES ($1, $2, $3)"

	now := time.Now()

	_, err = r.db.Master().Exec(query, user.Id, friend.Id, now)

	if err != nil {
		return app_error.NewInternalServerError(err)
	}

	return nil
}

func (r *FriendRepositoryInstance) RemoveFriend(user *user.User, friend *user.User) *app_error.AppError {
	rows, err := r.db.Slave().Query("SELECT EXISTS(SELECT 1 FROM friends WHERE user_id=$1 AND friend_id=$2)", user.Id, friend.Id)
	if err != nil {
		return app_error.NewInternalServerError(err)
	}

	var exists bool

	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&exists); err != nil {
			return app_error.NewInternalServerError(err)
		}

		if !exists {
			return nil
		}
	}

	query := "DELETE FROM friends WHERE user_id=$1 AND friend_id=$2"

	_, err = r.db.Master().Exec(query, user.Id, friend.Id)

	if err != nil {
		return app_error.NewInternalServerError(err)
	}

	return nil
}
