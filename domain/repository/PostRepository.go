package repository

import (
	"forum/domain/entity"

	"github.com/google/uuid"
)

type PostRepository interface {
	Create(post *entity.Post) error
	GetByID(postID uuid.UUID) (*entity.Post, error)
	GetbyuserId(Userid uuid.UUID)([]*entity.PostWithDetails)
	GetAll() ([]*entity.Post, error)
	GetByCategory(categoryID uuid.UUID) ([]*entity.Post, error)
	Update(post *entity.Post) error
	Delete(postID uuid.UUID) error
	GetWithDetails(postID uuid.UUID) (*entity.PostWithDetails, error)
	GetFiltered(filter entity.PostFilter) ([]*entity.Post, error)
}
