package entity

import "github.com/google/uuid"

type PostCategory struct {
	PostID     uuid.UUID `json:"post_id" db:"post_id"`
	CategoryID uuid.UUID `json:"category_id" db:"category_id"`
}
