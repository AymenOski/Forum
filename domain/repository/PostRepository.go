package repository

import (
	"forum/domain/entity"

	"github.com/google/uuid"
)
// methods
type PostRepository interface {
	Create(post *entity.Post) error
	GetByID(postID uuid.UUID) (*entity.Post, error)
	GetAll() ([]*entity.Post, error)
	GetByCategory(categoryID uuid.UUID) ([]*entity.Post, error)
	Update(post *entity.Post) error
	Delete(postID uuid.UUID) error
	GetFiltered(filter entity.PostFilter) ([]*entity.Post, error)
	// Me
	// category
	GetWithDetails(postID uuid.UUID) (*entity.PostWithDetails, error)
	// my posts
	GetbyuserId(Userid uuid.UUID) ([]*entity.Post, error)
	// liked posts
	//GetLikedPostsByUser(userID uuid.UUID) ([]*entity.PostWithDetails, error)
}
