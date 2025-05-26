package entity

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	UserID        uuid.UUID  `json:"user_id"`
	Name          string     `json:"name"`
	Email         string     `json:"email"`
	PasswordHash  string     `json:"-"` // Don't expose password hash
	SessionToken  *string    `json:"-"` // Don't expose session token
	SessionExpiry *time.Time `json:"-"` // Don't expose session expiry
	CreatedAt     time.Time  `json:"created_at"`
}
