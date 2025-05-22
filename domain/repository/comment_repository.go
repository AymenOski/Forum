package repository

import (
	"forum/domain/entity"

	"github.com/google/uuid"
)

type CommentRepository interface {
	Create(comment *entity.Comment) error
	GetByID(commentID int) (*entity.Comment, error)
	
	GetByPostID(postID int)([]*entity.Comment, error)
}