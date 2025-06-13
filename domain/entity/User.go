package entity

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `json:"id" db:"id"`
	UserName     string    `json:"user_name" db:"user_name"`
	Email        string    `json:"email" db:"email"`
	PasswordHash string    `json:"-" db:"password_hash"` // Don't expose password hash
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}
