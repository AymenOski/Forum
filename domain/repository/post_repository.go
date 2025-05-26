package repository

import (
	"forum/domain/entity"

	"github.com/google/uuid"
)


type PostRepository interface {
	Create(post *entity.Post) error

	GetAll() ([]*entity.Post, error)
	GetByUserID(userID *uuid.UUID) ([]*entity.Post, error)
	GetLikedPost() ([]*entity.Post, error) // 
	GetByCategory(categoryID []uint8) ([]*entity.Post, error)
}
