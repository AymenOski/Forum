package entity

import (
	"time"

	"github.com/google/uuid"
)

type Post struct {
	UserID        uuid.UUID `json:"user_id"`
	PostID        int `json:"post_id"`
	Authorname    string    `json:"username"`
	Content       string    `json:"content"`
	LikesCount    int       `json:"likes_count"`
	DislikesCount int       `json:"dislikes_count"`
	CreatedAt     time.Time `json:"created_at"`
}
