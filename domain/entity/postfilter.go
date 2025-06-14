package entity

import "github.com/google/uuid"

type PostFilter struct {
	CategoryIDs []uuid.UUID
	AuthorID    *uuid.UUID
	MyPosts     bool
	LikedPosts  bool
}
