package repository

import (
	"forum/domain/entity"

	"github.com/google/uuid"
)

type PostAggregateRepository interface {
	CreatePostWithCategories(post *entity.Post, categoryIDs []*uuid.UUID) error
	GetPostWithAllDetails(postID uuid.UUID) (*entity.PostWithDetails, error)
	GetFeedForUser() ([]*entity.PostWithDetails, error)
	GetPostsWithDetailsByUser(userID uuid.UUID) ([]*entity.PostWithDetails, error)
	GetFilteredPostsWithDetails(filter entity.PostFilter) ([]*entity.PostWithDetails, error)
}
