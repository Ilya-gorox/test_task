package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Firstname string    `json:"firstname" db:"firstname"`
	Lastname  string    `json:"lastname" db:"lastname"`
	Email     string    `json:"email" db:"email"`
	Age       uint      `json:"age" db:"age"`
	Created   time.Time `json:"created" db:"created"`
}

type CreateUserInput struct {
	Firstname string `json:"firstname" binding:"required"`
	Lastname  string `json:"lastname" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Age       uint   `json:"age" binding:"required,min=1"`
}

type UpdateUserInput struct {
	Firstname *string `json:"firstname"`
	Lastname  *string `json:"lastname"`
	Email     *string `json:"email" binding:"omitempty,email"`
	Age       *uint   `json:"age" binding:"omitempty,min=1"`
}
