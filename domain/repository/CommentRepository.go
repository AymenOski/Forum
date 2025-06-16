package repository

import (
	"forum/domain/entity"

	"github.com/google/uuid"
)

type CommentRepository interface {
	Create(comment *entity.Comment) error
	GetByID(commentID uuid.UUID) (*entity.Comment, error)
	GetByPostID(postID uuid.UUID) ([]entity.Comment, error)
	Delete(commentID uuid.UUID) error
	GetCountByPostID(postID uuid.UUID) (int, error)
	GetCountByUserID(userID uuid.UUID) (int, error)
	GetWithDetails(commentID uuid.UUID) (*entity.CommentWithDetails, error)
	GetByPostIDWithDetails(postID uuid.UUID) ([]entity.CommentWithDetails, error)
}
