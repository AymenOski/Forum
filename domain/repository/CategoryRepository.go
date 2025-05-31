package repository

import (
	"forum/domain/entity"

	"github.com/google/uuid"
)

type CategoryRepository interface {
	Create(category *entity.Category) error
	GetByID(categoryID uuid.UUID) (*entity.Category, error)
	GetByName(name string) (*entity.Category, error)
	GetAll() ([]*entity.Category, error)
	// Update(category *entity.Category) error
	// Delete(categoryID uuid.UUID) error
	CheckNameExists(name string) (bool, error)
	// GetWithPostCount() ([]*entity.Category, error)
}
