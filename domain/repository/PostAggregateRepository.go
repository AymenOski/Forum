package repository

import (
	"forum/domain/entity"

	"github.com/google/uuid"
)

type PostAggregateRepository interface {
	CreatePostWithCategories(post *entity.Post, categoryIDs []*uuid.UUID) error
	GetFeedForUser() ([]*entity.PostWithDetails, error)
	// Me
	GetPostWithAllDetails(postID uuid.UUID) (*entity.PostWithDetails, error)
	GetPostsWithDetailsByUser(userID uuid.UUID) ([]*entity.PostWithDetails, error)
	GetLikedPostsByUser(userID uuid.UUID) ([]*entity.PostWithDetails, error)
}
