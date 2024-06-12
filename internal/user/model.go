package user

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	Id        uuid.UUID `json:"id" db:"id"`
	Name      string    `json:"name" binding:"required,alphanum"`
	Surname   string    `json:"surname" binding:"required,alphanum"`
	Email     string    `json:"email" binding:"required,email"`
	Birthday  time.Time `json:"birthday" binding:"required"`
	Biography string    `json:"biography" binding:"alphanum" faker:"lang=rus, sentence, "`
	City      string    `json:"city" binding:"alphanum" faker:"city, lang=rus"`
	Password  string    `json:"-" binding:"required,alphanum" faker:"word, lang=rus"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
