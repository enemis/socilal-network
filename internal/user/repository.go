package user

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"social-network-otus/internal/app_error"
	"social-network-otus/internal/database"
	"social-network-otus/internal/rest/response"
)

type UserRepository interface {
	GetUserByEmail(email string) (*User, *app_error.AppError)
	GetUserById(userId string) (*User, *app_error.AppError)
	CreateUser(user *User) (*uuid.UUID, *app_error.AppError)
	FindUsers(name, surname string) ([]*User, *app_error.AppError)
}

type UserRepositoryInstance struct {
	db       *database.DatabaseStack
	response *response.ResponseFactory
}

func NewUserRepository(databaseStack *database.DatabaseStack) *UserRepositoryInstance {
	return &UserRepositoryInstance{db: databaseStack}
}

func (r *UserRepositoryInstance) GetUserByEmail(email string) (*User, *app_error.AppError) {
	rows, err := r.db.Slave().Queryx("SELECT * FROM users WHERE email=$1", email)
	if err != nil {
		return nil, app_error.NewInternalServerError(err)
	}
	defer rows.Close()

	var user User
	for rows.Next() {
		err = rows.StructScan(&user)
		if err != nil {
			return nil, app_error.NewInternalServerError(err)
		}
	}

	return &user, nil
}

func (r *UserRepositoryInstance) GetUserById(userId string) (*User, *app_error.AppError) {
	var user User
	err := r.db.GetReadConnection().Get(&user, "SELECT * FROM users WHERE id=$1 LIMIT 1", userId)

	if err != nil {
		return nil, app_error.NewBadRequestFromError(errors.New("user not found"))
	}

	return &user, nil
}

func (r *UserRepositoryInstance) FindUsers(name, surname string) ([]*User, *app_error.AppError) {
	users := make([]*User, 100)
	query := "SELECT * FROM users WHERE "
	paramName := strings.ToLower(name) + "%"
	paramSurname := strings.ToLower(surname) + "%"
	limitPart := " ORDER BY id LIMIT 100;"

	var err error

	if len(name) > 1 && len(surname) > 1 {
		err = r.db.Slave().Select(&users, query+"(lower(name) LIKE $1 AND lower(surname) LIKE $2) OR (lower(surname) LIKE $3 AND lower(name) LIKE $4)"+limitPart, paramName, paramSurname, paramName, paramSurname)
	} else if len(name) > 0 {
		err = r.db.Slave().Select(&users, query+"lower(name) LIKE $1"+limitPart, paramName)
	} else if len(surname) > 0 {
		err = r.db.Slave().Select(&users, query+"lower(surname) LIKE $1"+limitPart, paramSurname)
	} else {
		err = r.db.Slave().Select(&users, query+limitPart)
	}

	if err != nil {
		return users, app_error.NewBadRequestFromError(errors.New("user not found"))
	}

	return users, nil
}

func (r *UserRepositoryInstance) CreateUser(user *User) (*uuid.UUID, *app_error.AppError) {
	rows, err := r.db.Slave().Query("SELECT EXISTS(SELECT 1 FROM users WHERE email=$1)", user.Email)
	if err != nil {
		return nil, app_error.NewInternalServerError(err)
	}

	var exists bool

	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&exists); err != nil {
			return nil, app_error.NewInternalServerError(err)
		}

		if exists {
			return nil, app_error.NewBadRequestFromError(errors.New(fmt.Sprintf("User with email %s already registered", user.Email)))
		}
	}

	query := "INSERT INTO users (name, surname, email, birthday, biography, city, password, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id"

	now := time.Now()
	var userId uuid.UUID
	err = r.db.GetWriteConnection().QueryRow(query, user.Name, user.Surname, user.Email, user.Birthday, user.Biography, user.City, user.Password, now, now).Scan(&userId)

	if err != nil {
		return nil, app_error.NewInternalServerError(err)
	}

	return &userId, nil
}
