package usecase

import (
	"errors"
	"time"

	"forum/domain/entity"
	"forum/domain/repository"

	"github.com/google/uuid"
)

type PostService struct {
	postRepo         repository.PostRepository
	userRepo         repository.UserRepository
	categoryRepo     repository.CategoryRepository
	postCategoryRepo repository.PostCategoryRepository
}

func NewPostService(postRepo repository.PostRepository, userRepo repository.UserRepository,
	categoryRepo repository.CategoryRepository, postCategoryRepo repository.PostCategoryRepository) *PostService {
	return &PostService{
		postRepo:         postRepo,
		userRepo:         userRepo,
		categoryRepo:     categoryRepo,
		postCategoryRepo: postCategoryRepo,
	}
}

func (ps *PostService) CreatePost(userID *uuid.UUID, authorName string, content string, categoryIDs []*uuid.UUID) (*entity.Post, error) {
	user, err := ps.userRepo.GetByID(*userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}
	if content == "" {
		return nil, errors.New("post content cannot be emtpy")
	}
	if len(categoryIDs) <= 0 {
		return nil, errors.New("you have to select one category at least")
	}
	if len(content) > 5000 {
		return nil, errors.New("post content too long (max: 5000 character)")
	}

	for _, categoryID := range categoryIDs {
		category, err := ps.categoryRepo.GetByID(*categoryID)
		if err != nil {
			return nil, err
		}
		if category == nil {
			return nil, errors.New("one or more categories does not exist")
		}
	}
	post := &entity.Post{
		UserID:    *userID,
		Content:   content,
		CreatedAt: time.Now(),
	}
	err = ps.postRepo.Create(post)
	if err != nil {
		return nil, err
	}
	// Associate the categories to the post
	err = ps.postCategoryRepo.AddCategoriesToPost(post.ID, categoryIDs)
	if err != nil {
		err := ps.postRepo.Delete(post.ID)
		if err != nil {
			return nil, err
		}
		return nil, err
	}
	return post, nil
}
