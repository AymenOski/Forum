package repository

import (
	"forum/domain/entity"

	"github.com/google/uuid"
)

type PostRepository interface {
	Create(post *entity.Post) error
	GetByID(postID uuid.UUID) (*entity.Post, error)
	GetAll() ([]*entity.Post, error)
	GetByUserID(userID uuid.UUID) ([]*entity.Post, error)
	GetWithPagination(limit, offset int) ([]*entity.Post, error)
	GetByCategory(categoryID uuid.UUID) ([]*entity.Post, error)
	GetByCategoryWithPagination(categoryID uuid.UUID, limit, offset int) ([]*entity.Post, error)
	GetMostLiked(limit int) ([]*entity.Post, error)
	GetRecent(limit int) ([]*entity.Post, error)
	Update(post *entity.Post) error
	Delete(postID uuid.UUID) error
	Search(query string) ([]*entity.Post, error)
	GetWithDetails(postID uuid.UUID) (*entity.PostWithDetails, error)
	GetAllWithDetails() ([]*entity.PostWithDetails, error)
}
