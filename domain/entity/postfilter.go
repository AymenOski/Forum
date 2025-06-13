package entity

import "github.com/google/uuid"

type PostFilter struct {
	CategoryID *uuid.UUID
	AuthorID   *uuid.UUID
}
