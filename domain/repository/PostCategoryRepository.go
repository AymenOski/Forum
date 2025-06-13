package repository

import (
	"forum/domain/entity"

	"github.com/google/uuid"
)

type PostCategoryRepository interface {
	Create(postCategory *entity.PostCategory) error
	Delete(postID, categoryID uuid.UUID) error
	DeleteByPostID(postID uuid.UUID) error
	DeleteByCategoryID(categoryID uuid.UUID) error
	GetCategoriesByPostID(postID uuid.UUID) ([]*entity.Category, error)
	GetPostsByCategoryID(categoryID uuid.UUID) ([]*entity.Post, error)
	CheckAssociationExists(postID, categoryID uuid.UUID) (bool, error)
	GetAllAssociations() ([]*entity.PostCategory, error)
}