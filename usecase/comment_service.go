package usecase

import (
	"errors"
	"forum/domain/entity"
	"forum/domain/repository"
	"strings"
	"time"

	"github.com/google/uuid"
)

type CommentService struct {
	userRepo    repository.UserRepository
	commentRepo repository.CommentRepository
	postRepo    repository.PostRepository
}

func NewCommentService(userRepo repository.UserRepository, commentRepo repository.CommentRepository,
	postRepo repository.PostRepository) *CommentService {
	return &CommentService{
		userRepo:    userRepo,
		commentRepo: commentRepo,
		postRepo:    postRepo,
	}
}

func (cs *CommentService) CreateComment(postID *uuid.UUID, userID *uuid.UUID, content string) (*entity.Comment, error) {
	content = strings.TrimSpace(content)
	if len(content) > 250 {
		return nil, errors.New("comment length excceds 250 characters")
	} else if content == "" {
		return nil, errors.New("comment should have at least 1 character")
	}

	_, err := cs.postRepo.GetByID(*postID)
	if err != nil {
		return nil, errors.New("post not found")
	}
	_, err = cs.userRepo.GetByID(*userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	comment := &entity.Comment{
		Content:   content,
		UserID:    *userID,
		PostID:    *postID,
		CreatedAt: time.Now(),
	}

	err = cs.commentRepo.Create(comment)
	if err != nil {
		return nil, err
	}

	return comment, nil
}
