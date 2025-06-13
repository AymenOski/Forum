package entity

import (
	"time"

	"github.com/google/uuid"
)

type CommentReaction struct {
	ID        uuid.UUID `json:"id" db:"id"`
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	CommentID uuid.UUID `json:"comment_id" db:"comment_id"`
	Reaction  bool      `json:"reaction" db:"reaction"` // true for like, false for dislike
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
