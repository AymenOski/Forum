package repository

import (
	"forum/domain/entity"

	"github.com/google/uuid"
)

type LikeRepository interface {
	AddParentReaction(postID int, userID *uuid.UUID, isLike bool) error
	RemoveParentReaction(postID int, userID *uuid.UUID) error
	GetParentReaction(postID int, userID *uuid.UUID) (*entity.LikeDislike, error) 
	CountParentLikes(postID int)(int, error)
	CountParentDislikes(postID int)(int, error)
}