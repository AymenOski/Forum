package repository

import "forum/domain/entity"

type PostCategoryRepository interface {
	AddCategoriesToPost(postID int, categoryID []uint8) error
	GetCategoriesByPostID(postID int) ([]*entity.Category, error)
	GetPostsByCategory(categoryID uint8) ([]*entity.Post, error)
}
