package repository

import "forum/domain/entity"

type CategoryRepository interface {
	GetByID(categoryID uint8) (*entity.Category, error)
	GetAll() ([]*entity.Category, error)
}
