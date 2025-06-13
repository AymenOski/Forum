package repository

import (
	"forum/domain/entity"

	"github.com/google/uuid"
)

type PostAggregateRepository interface {
	CreatePostWithCategories(post *entity.Post, categoryIDs []*uuid.UUID) error
	GetPostWithAllDetails(postID uuid.UUID) (*entity.PostWithDetails, error)
	GetFeedForUser() ([]*entity.PostWithDetails, error)
}
