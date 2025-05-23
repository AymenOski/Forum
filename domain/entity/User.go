package entity

import "time"

// User represents a user in the system
type User struct {
	UserID        string     `json:"user_id"`
	Name          string     `json:"name"`
	Email         string     `json:"email"`
	PasswordHash  string     `json:"-"`
	SessionToken  *string    `json:"-"`
	SessionExpiry *time.Time `json:"-"`
}
