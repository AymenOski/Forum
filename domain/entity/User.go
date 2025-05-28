package entity

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	UserID        uuid.UUID  `json:"user_id"`
	Name          string     `json:"name"`
	Email         string     `json:"email"`
	PasswordHash  string     `json:"-"`
	SessionToken  *string    `json:"-"`
	SessionExpiry *time.Time `json:"-"`
	CreatedAt     time.Time  `json:"created_at"`
}
