package repository

import (
	"forum/domain/entity"

	"github.com/google/uuid"
)

type PostReactionRepository interface {
	Create(reaction *entity.PostReaction) error
	GetByID(reactionID uuid.UUID) (*entity.PostReaction, error)
	GetByUserAndPost(userID, postID uuid.UUID) (*entity.PostReaction, error)
	GetByPostID(postID uuid.UUID) ([]*entity.PostReaction, error)
	GetByUserID(userID uuid.UUID) ([]*entity.PostReaction, error)
	Update(reaction *entity.PostReaction) error
	Delete(reactionID uuid.UUID) error
	DeleteByUserAndPost(userID, postID uuid.UUID) error
	GetLikeCountByPostID(postID uuid.UUID) (int, error)
	GetDislikeCountByPostID(postID uuid.UUID) (int, error)
	GetReactionCountsByPostID(postID uuid.UUID) (likes int, dislikes int, err error)
	HasUserReacted(userID, postID uuid.UUID) (bool, *bool, error) // exists, reaction_value, error
}
