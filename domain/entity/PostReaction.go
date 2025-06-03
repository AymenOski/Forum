package entity

import (
	"time"

	"github.com/google/uuid"
)

type PostReaction struct {
	ID        uuid.UUID `json:"id" db:"id"`
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	PostID    uuid.UUID `json:"post_id" db:"post_id"`
	Reaction  bool      `json:"reaction" db:"reaction"` // true for like, false for dislike
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
