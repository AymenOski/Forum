package repository

import (
	"time"

	"forum/domain/entity"

	"github.com/google/uuid"
)

type PostAggregateRepository interface {
	CreatePostWithCategories(post *entity.Post, categoryIDs []*uuid.UUID) error
	GetPostWithAllDetails(postID uuid.UUID) (*entity.PostWithDetails, error)
	GetFeedForUser(userID uuid.UUID, limit, offset int) ([]*entity.PostWithDetails, error)
	GetTrendingPosts(since time.Time, limit int) ([]*entity.PostWithDetails, error)
}
