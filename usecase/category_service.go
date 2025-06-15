package usecase

import (
	"errors"
	"strings"

	"forum/domain/entity"
	"forum/domain/repository"
)

type CategoryService struct {
	categoryRepo     repository.CategoryRepository
	postCategoryRepo repository.PostCategoryRepository
	sessionRepo      repository.UserSessionRepository
	userRepo         repository.UserRepository
}

func NewCategoryService(
	categoryRepo repository.CategoryRepository,
	postCategoryRepo repository.PostCategoryRepository,
	sessionRepo repository.UserSessionRepository,
	userRepo repository.UserRepository,
) *CategoryService {
	return &CategoryService{
		categoryRepo:     categoryRepo,
		postCategoryRepo: postCategoryRepo,
		sessionRepo:      sessionRepo,
		userRepo:         userRepo,
	}
}

func (cs *CategoryService) GetCategoryByName(name string) (*entity.Category, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, errors.New("category name cannot be empty")
	}

	category, err := cs.categoryRepo.GetByName(name)
	if err != nil {
		return nil, err
	}

	return category, nil
}

func (cs *CategoryService) GetAllCategories() ([]*entity.Category, error) {
	return cs.categoryRepo.GetAll()
}