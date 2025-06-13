package entity

import (
	"time"

	"github.com/google/uuid"
)

type Comment struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Content   string    `json:"content" db:"content"`
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	PostID    uuid.UUID `json:"post_id" db:"post_id"`
	CreatedAt time.Time `json:"createdat" db:"createdat"` // Note: schema shows 'createdat' not 'created_at'
}
