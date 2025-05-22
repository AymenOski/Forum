package repository

import (
	"forum/domain/entity"

	"github.com/google/uuid"
)

type LikeRepository interface {
	AddPostReaction(postID int, userID *uuid.UUID, isLike bool) error
	RemovePostReaction(postID int, userID *uuid.UUID) error
	GetPostReaction(postID int, userID *uuid.UUID) (*entity.PostLike, error)

	CountPostLikes(postID int)(int, error)
	CountPostDislikes(postID int)(int, error)
}